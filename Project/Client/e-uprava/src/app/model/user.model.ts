export class User {
  username: string;
  password: string
  fname: string;
  lname: string;
  usertype: string;

  constructor( email: string, password: string, first_name: string, last_name: string, role: string) {
    this.username = email;
    this.password = password;
    this.fname = first_name;
    this.lname = last_name;
    this.usertype = role;
  }
}
