const uploadForm = document.querySelector(".upload");
const fileInput = document.querySelector(".file-input")

uploadForm.addEventListener("click", () => {
    fileInput.click();
});

fileInput.onchange = function () {
    uploadFile();
}

function uploadFile() {
    const xhr = new XMLHttpRequest();
    xhr.open("put", "/api/upload");
    xhr.responseType = `json`;
    const data = new FormData();

    var files = document.getElementById('file').files;
    for (var i = 0; i < files.length; i++) {
        data.append(files[i].name, files[i]);
    }

    xhr.send(data);
}



window.onload = function () {
    function getFileList() {
        const xhr = new XMLHttpRequest();
        xhr.open("get", "/api/filelist");
        xhr.setRequestHeader('Content-Type', 'application/json');

        xhr.responseType = `json`;
        xhr.onload = function () {
            if (xhr.status === 200) {
                const resp = xhr.response;
                const data = resp.data
                if (data.files?.length) {
                    Array.from(data.files).forEach(file => {
                        const li = document.createElement("li");
                        li.innerHTML = `<a download href="/uploads/${data.userid}/${file}">${file}</a>`;
                        list.appendChild(li);
                    });
                }
            }
        };
        xhr.send();
    }
    getFileList();
}
