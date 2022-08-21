window.onload = function() {
    async function getFileList() {
        var xhr = new XMLHttpRequest();
        xhr.open("get", "/api/filelist");
        xhr.setRequestHeader('Content-Type', 'application/json');
        
        xhr.responseType = `json`;
        xhr.onload = function() {
            var status = xhr.status
            if (status === 200) {
                const data = xhr.response;
                const userid = data.userid;
                Array.from(data.files).forEach(file => {
                    var li = document.createElement("li");
                    li.innerHTML = `<a download href="/uploads/${userid}/${file}">${file}</a>`;
                    list.appendChild(li);
                });
            } else {
                alert(status.toString() + ` ` + xhr.statusText);
                return;
            }
        };
        xhr.send();
        


        /*let req = await fetch(`/api/filelist`).then(r => r.json());

        let list = document.getElementById("list");

        Object.keys(req).forEach(element => {
            Array.from(req[element]).forEach(file => {
                var li = document.createElement("li");
                li.innerHTML = `<a download href="/uploads/${userid}/${file}">${file}</a>`;
                list.appendChild(li);
            });
        });*/
    }
    getFileList();
}
