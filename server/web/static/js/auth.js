window.onload = function () {
  let authApiPath = undefined;
  function handleAuth(event, path) {
    event.preventDefault();

    const data = new FormData(event.target);
    const value = Object.fromEntries(data.entries());

    var xhr = new XMLHttpRequest();
    xhr.open("post", authApiPath);
    xhr.setRequestHeader('Content-Type', 'application/json');

    xhr.responseType = `json`;
    xhr.onload = function () {
      var status = xhr.status;
      if (status === 200) {
        const resp = xhr.response;
        const data = resp.data;
        document.cookie = `login=${data.login}; SameSite=None; Secure`;
        document.cookie = `token=${data.token}; expires=${data.expires}; SameSite=None; Secure`;
        window.location.replace("/profile");
      } else {
        alert(status.toString() + ` ` + xhr.statusText);
        return;
      }
    };
    xhr.send(JSON.stringify(value));
  }

  if (document.getElementById('register_form') != null) {
    authApiPath = "/api/register";
    const register = document.querySelector('#register_form');
    register.addEventListener('submit', handleAuth);
    return
  }

  if (document.getElementById('login_form') != null) {
    authApiPath = "/api/login";
    const login = document.querySelector('#login_form');
    login.addEventListener('submit', handleAuth);
    return
  }
}