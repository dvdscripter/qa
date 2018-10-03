import * as jwt from "jsonwebtoken";

const isLogged = () => {
  let token = localStorage.getItem("token");
  if (!token) return false;
  let decoded = jwt.decode(token);
  if (decoded) return decoded.email;
  return false;
};

export default isLogged;
