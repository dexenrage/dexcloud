var authLinkParams = `href="/login"`;
var authLinkIcon = `/static/img/login.svg`
var authLinkWord = `Войти`;


function checkAuth() {
    const xhr = new XMLHttpRequest();
    xhr.open("get", `/api/checkauth`)
    xhr.setRequestHeader('Content-Type', 'application/json');

    xhr.responseType = `json`;

    var xhrStatus = undefined;
    xhr.onload = function () {
        xhrStatus = xhr.status
        if (xhrStatus == 200) {
            authLinkParams = ``;
            authLinkIcon = `/static/img/logout.svg`;
            authLinkWord = `Выйти`;
            addHeader();

            const element = document.getElementById("header-auth-link");
            element.onclick = function () {
                document.cookie = `token` + '=;expires=Thu, 01 Jan 1970 00:00:01 GMT; SameSite=None; Secure';
                document.cookie = `login` + '=;expires=Thu, 01 Jan 1970 00:00:01 GMT; SameSite=None; Secure';
                window.location.replace("/");
            };

            return

        } else {
            addHeader();
            return
        }
    }
    xhr.send();
}

function addHeader() {

    var headerAuthLink = `
    <a class="header-auth-link" id="header-auth-link" ${authLinkParams}>
        <img alt="${authLinkWord}" src="${authLinkIcon}">
        <b>${authLinkWord}</b>
    </a>`;

    if (window.location.pathname === `/login` || window.location.pathname === `/register`) {
        headerAuthLink = ``;
    }

    const headerText = `
        <div class="header">
            <link rel="stylesheet" href="/static/css/header.css" type="text/css">
                <b class="header-title">DexCloud</b>
            ${headerAuthLink}
        </div>`;

    document.body.insertAdjacentHTML('afterbegin', headerText);
}

checkAuth();
