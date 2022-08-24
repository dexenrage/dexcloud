function addHeader() {
    var headerAuthLink = `
        <a class="header-auth-link" href="/login">
            <img alt="Вход" src="/static/img/account.svg">
            <b>Вход</b>
        </a>`;

    if (window.location.pathname === `/login` || window.location.pathname === `/register`) {
        headerAuthLink = ``;
    }

    var headerText = `
        <div class="header">
            <link rel="stylesheet" href="/static/css/header.css" type="text/css">
            
                <b class="header-title">DexCloud</b>
            
            ${headerAuthLink}
        </div>`;

    document.body.insertAdjacentHTML('afterbegin', headerText);
}
addHeader();
