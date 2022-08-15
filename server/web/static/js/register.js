window.onload = function() {
    function handleSubmit(event) {
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
            document.cookie = `token=${data.token}; SameSite=None; Secure`;
          } else {
            alert(status.toString() + ` ` + xhr.statusText)
            return
          }
        };
        xhr.send(JSON.stringify(value))
      }
      const form = document.querySelector('form');
      form.addEventListener('submit', handleSubmit);
}