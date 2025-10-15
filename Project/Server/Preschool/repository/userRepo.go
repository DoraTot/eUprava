package repository

//
//import (
//	"crypto/rsa"
//	"database/sql"
//	"encoding/json"
//	"fmt"
//	"github.com/golang-jwt/jwt/v5"
//	"io"
//	"main.go/config"
//	"main.go/model"
//	"net/http"
//	"strings"
//)
//
//type UserRepo struct {
//	DB *sql.DB
//}
//
//func NewUserRepo(db *sql.DB) *UserRepo {
//	return &UserRepo{DB: db}
//}
//
//func GetAuth0PublicKey(kid string) (*rsa.PublicKey, error) {
//	resp, err := http.Get("https://" + config.Auth0Domain + "/.well-known/jwks.json")
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	body, _ := io.ReadAll(resp.Body)
//	var jwksData model.Jwks
//	if err := json.Unmarshal(body, &jwksData); err != nil {
//		return nil, err
//	}
//
//	for _, key := range jwksData.Keys {
//		if key.Kid == kid {
//			return jwt.ParseRSAPublicKeyFromPEM([]byte(convertX5CToPEM(key.X5c[0])))
//		}
//	}
//	return nil, fmt.Errorf("public key not found for kid %s", kid)
//}
//
//func convertX5CToPEM(cert string) string {
//	return "-----BEGIN CERTIFICATE-----\n" + cert + "\n-----END CERTIFICATE-----"
//}
//
//func GetParentsFromAuth0() ([]model.ParentUser, error) {
//	token, err := getManagementToken()
//	if err != nil {
//		return nil, err
//	}
//
//	fmt.Println("🔑 Management Token:", token)
//	fmt.Printf("Token length: %d\n", len(token))
//
//	roleID := "rol_Hadxe5arTWn2QgVj"
//	domain := config.Auth0Domain
//
//	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/api/v2/roles/%s/users", domain, roleID), nil)
//	req.Header.Set("Authorization", "Bearer "+token)
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("Auth0 API error: %s", string(body))
//	}
//
//	//var result struct {
//	//	Users []model.ParentUser `json:"users"`
//	//}
//	var result []model.ParentUser
//	if err := json.Unmarshal(body, &result); err != nil {
//		return nil, err
//	}
//
//	//return result.Users, nil
//	return result, nil
//
//	//body, err := io.ReadAll(resp.Body)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//if resp.StatusCode != http.StatusOK {
//	//	return nil, fmt.Errorf("Auth0 API error: %s", string(body))
//	//}
//
//	//body, _ := io.ReadAll(resp.Body)
//	//fmt.Println("Auth0 response:", string(body))
//	//log.Println("Auth0 response:", string(body))
//
//	//if resp.StatusCode != http.StatusOK {
//	//	return nil, fmt.Errorf("Auth0 API error: %s", resp.Body)
//	//}
//
//	//var users []model.ParentUser
//	//if err := json.Unmarshal(body, &users); err != nil {
//	//	return nil, err
//	//}
//	//return users, nil
//
//	//var result struct {
//	//	Users []model.ParentUser `json:"users"`
//	//}
//	//if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
//	//	return nil, err
//	//}
//	//return result.Users, nil
//}
//
//func getManagementToken() (string, error) {
//	payload := strings.NewReader(fmt.Sprintf(`{
//        "client_id":"%s",
//        "client_secret":"%s",
//        "audience":"%s",
//        "grant_type":"client_credentials"
//    }`, config.Auth0Client, config.Auth0Secret, config.Audience))
//
//	resp, err := http.Post("https://"+config.Auth0Domain+"/oauth/token", "application/json", payload)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	var result struct {
//		AccessToken string `json:"access_token"`
//	}
//	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
//		return "", err
//	}
//
//	return result.AccessToken, nil
//}
