window.onload = function() {
    async function getFileList() {
        let req = await fetch(`/api/filelist`).then(r => r.json());

        let list = document.getElementById("list");

        Object.keys(req).forEach(element => {
            Array.from(req[element]).forEach(file => {
                var li = document.createElement("li");
                li.innerHTML = `<a download href="/uploads/1/${file}">${file}</a>`;
                list.appendChild(li);
            });
        });
    }
    getFileList();
}