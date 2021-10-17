package ssojwt

const wait = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <title>Please Wait</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
      body {
        background-color: white;
        color: #5d31ff;
        font-family: -apple-system, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Open Sans",
          "Helvetica Neue", sans-serif;
        font-size: 16px;
        text-align: center;
      }
    </style>
  </head>
  <body>
    <div id="wrap">
      <h3>Please Wait</h3>
    </div>
    <script type="application/javascript">
      (() => {
        if (window.opener) {
          const rawData = "{{.LoginResponse}}";
          const data = JSON.parse(rawData);
          console.log(data);
          window.opener.postMessage(data, "{{.OriginUrl}}");
        } else {
          console.log("idk");
        }
      })();
    </script>
  </body>
</html>
`
