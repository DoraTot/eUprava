package repository

import (
	"database/sql"
	"log"
	"main.go/model"
)

type ChildRepo struct {
	DB *sql.DB
}

func NewChildrenRepo(db *sql.DB) *ChildRepo {
	return &ChildRepo{DB: db}
}

func (r *ChildRepo) EnsureTableExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS children (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		parent_id VARCHAR(255) NOT NULL,
		enrolled VARCHAR(20) NOT NULL DEFAULT 'PENDING',
		exam_done BOOLEAN NOT NULL DEFAULT FALSE
	);`
	_, err := r.DB.Exec(query)
	if err != nil {
		log.Println("Failed to create children table:", err)
	}
	return err
}

func (r *ChildRepo) Create(name string, parentId string, enrolled string) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO children (name, parent_id, enrolled, exam_done)
                           VALUES (?, ?, ?, false)`, name, parentId, enrolled)
	if err != nil {
		log.Println("Failed to insert child:", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ChildRepo) GetAllChildren() ([]model.Child, error) {
	rows, err := r.DB.Query(`SELECT id, name, parent_id, exam_done, enrolled FROM children`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []model.Child
	for rows.Next() {
		var c model.Child
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentId, &c.ExamDone, &c.Enrolled); err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, nil
}

func (r *ChildRepo) GetChildByParentID(parentID string) ([]model.Child, error) {
	rows, err := r.DB.Query(`SELECT id, name, parent_id, exam_done, enrolled 
                             FROM children WHERE parent_id = ?`, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []model.Child
	for rows.Next() {
		var c model.Child
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentId, &c.ExamDone, &c.Enrolled); err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, nil
}

func (r *ChildRepo) GetChildrenWithoutExam() ([]model.Child, error) {
	rows, err := r.DB.Query(`SELECT id, name, parent_id, exam_done, enrolled 
                             FROM children WHERE exam_done = false`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []model.Child
	for rows.Next() {
		var c model.Child
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentId, &c.ExamDone, &c.Enrolled); err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, nil
}

func (r *ChildRepo) UpdateExamStatus(parentId string, childName string, examDone bool) (int64, error) {
	res, err := r.DB.Exec(`UPDATE children SET exam_done = ? WHERE parent_id = ? AND name = ?`,
		examDone, parentId, childName)
	if err != nil {
		log.Println("Failed to update exam:", err)
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (r *ChildRepo) UpdateEnrollmentStatus(parentId string, childName string, enrolled string) (int64, error) {
	res, err := r.DB.Exec(`UPDATE children SET enrolled = ? WHERE parent_id = ? AND name = ?`,
		enrolled, parentId, childName)
	if err != nil {
		log.Println("Failed to update enrollment:", err)
		return 0, err
	}

	return res.RowsAffected()
}

func (r *ChildRepo) GetChildByParentAndName(parentID string, childName string) (*model.Child, error) {

	query := `
		SELECT id, name, parent_id, exam_done, enrolled
		FROM children
		WHERE parent_id = ? AND name = ?
		LIMIT 1
	`

	row := r.DB.QueryRow(query, parentID, childName)

	var child model.Child

	err := row.Scan(
		&child.ID,
		&child.Name,
		&child.ParentId,
		&child.ExamDone,
		&child.Enrolled,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Println("Failed to fetch child:", err)
		return nil, err
	}

	return &child, nil
}
