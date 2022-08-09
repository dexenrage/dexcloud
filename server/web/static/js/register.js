window.onload = function() {
    function handleSubmit(event) {
        event.preventDefault();
      
        const data = new FormData(event.target);
        const value = Object.fromEntries(data.entries());

        var xhr = new XMLHttpRequest();
        xhr.open("post", "/api/register");
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.send(JSON.stringify(value))
      }
      const form = document.querySelector('form');
      form.addEventListener('submit', handleSubmit);
}