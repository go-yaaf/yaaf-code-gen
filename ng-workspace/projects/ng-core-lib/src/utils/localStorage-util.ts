// Access token key in the local storage
const tokenKey = 'portalAccessToken';
const loginKey = 'portalLoginData';

export function getToken() {
  return localStorage.getItem(tokenKey);
}

export function setToken(token) {
  localStorage.setItem(tokenKey, token);
}

export function removeToken() {
  localStorage.removeItem(tokenKey);
  localStorage.removeItem(loginKey);
}