window.onload = function () {
    function getFileList() {
        var xhr = new XMLHttpRequest();
        xhr.open("get", "/api/filelist");
        xhr.setRequestHeader('Content-Type', 'application/json');

        xhr.responseType = `json`;
        xhr.onload = function () {
            var status = xhr.status
            if (status === 200) {
                const resp = xhr.response;
                const data = resp.data
                if (data.files?.length) {
                    Array.from(data.files).forEach(file => {
                        var li = document.createElement("li");
                        li.innerHTML = `<a download href="/uploads/${data.userid}/${file}">${file}</a>`;
                        list.appendChild(li);
                    });
                }
            } else {
                alert(status.toString() + ` ` + xhr.statusText);
                return;
            }
        };
        xhr.send();
    }
    getFileList();
}
