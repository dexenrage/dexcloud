window.onload = function() {
    function handleRegister(event) {
        event.preventDefault();

        const data = new FormData(event.target);
        const value = Object.fromEntries(data.entries());

        var xhr = new XMLHttpRequest();
        xhr.open("post", "/api/register");
        xhr.setRequestHeader('Content-Type', 'application/json');

        xhr.responseType = `json`;
        xhr.onload = function() {
          var status = xhr.status;
          if (status === 200) {
            const data = xhr.response;
            document.cookie = `userid=${data.userid}; SameSite=None; Secure`;
            document.cookie = `token=${data.token}; SameSite=None; Secure`;
            window.location.replace("/profile");
          } else {
            alert(status.toString() + ` ` + xhr.statusText);
            return;
          }
        };
        xhr.send(JSON.stringify(value));
      }
      const register = document.querySelector('#register_form');
      register.addEventListener('submit', handleRegister);

      function handleLogin(event) {
        event.preventDefault();

        const data = new FormData(event.target);
        const value = Object.fromEntries(data.entries());

        var xhr = new XMLHttpRequest();
        xhr.open("post", "/api/login");
        xhr.setRequestHeader('Content-Type', 'application/json');

        xhr.responseType = `json`;
        xhr.onload = function() {
          var status = xhr.status;
          if (status === 200) {
            const data = xhr.response;
            document.cookie = `userid=${data.userid}; SameSite=None; Secure`;
            document.cookie = `token=${data.token}; SameSite=None; Secure`;
            window.location.replace("/profile");
          } else {
            alert(status.toString() + ` ` + xhr.statusText);
            return;
          }
        };
        xhr.send(JSON.stringify(value));
      }
      const login = document.querySelector('#login_form');
      login.addEventListener('submit', handleLogin);
}