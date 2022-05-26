# golang-sso-ui-jwt

Golang Library to use SSO UI as JWT login

## Installation

`go get github.com/ristekoss/golang-uisso-jwt`

## How to Use
see example folder


login handler
```go
// set config
config := ssojwt.MakeSSOConfig(time.Hour*168, time.Hour*720, "super secret access", "super secret refresh", "http://localhost:8080/login", "http://localhost:8080/")

// create middleware
authMiddleware := ssojwt.LoginCreator(config, nil)

// use handler func
http.HandleFunc("/login", authMiddleware)
```

on the frontend
```js
const loginHandler = async () => {
  const data = await popUpLogin();
  const e = document.getElementById("res");
  e.innerHTML = JSON.stringify(data, null, 2);
  localStorage.setItem("ssoui", JSON.stringify(data));
};

const popUpLogin = () => {
  const SSOWindow = window.open(
    new URL(
      // change service to what your backend url use and url encode it
      "https://sso.ui.ac.id/cas2/login?service=http%3A%2F%2Flocalhost%3A8080%2Flogin"
    ).toString(),
    "SSO UI Login",
    "left=50, top=50, width=480, height=480"
  );

  return new Promise(function (resolve, reject) {
    window.addEventListener(
      "message",
      (e) => {
        if (SSOWindow) {
          SSOWindow.close();
        }
        const data = e.data;
        resolve(data);
      },
      { once: true }
    );
  });
};

```

authenticated middleware
```go
// set config
config := ssojwt.MakeSSOConfig(time.Hour*168, time.Hour*720, "super secret access", "super secret refresh", "http://localhost:8080/login", "http://localhost:8080/")

// create middleware
middle := ssojwt.MakeAccessTokenMiddleware(config, "user")

// use handler func
auth := middle(handler)

http.Handle("/check", auth)
```

authenticated middleware
```go
// set config
config := ssojwt.MakeSSOConfig(time.Hour*168, time.Hour*720, "super secret access", "super secret refresh", "http://localhost:8080/login", "http://localhost:8080/")

// create middleware
middle := ssojwt.MakeRefreshTokenMiddleware(config)

// use handler func
http.Handle("/refresh", middle)
```
### To-dos

- create fiber handler
- create better readme.md
