function hideErrorMessage(id) {
    if(id===undefined) {
        id = "register-error"
    }
    document.getElementById(id).style.display ="none";
}

function showErrorMessage(id,msg) {
    if(id===undefined) {
        id = "register-error"
    }
    document.getElementById(id).innerHTML = msg;
    document.getElementById(id).style.display = "inline-block";
}

function isEmail(email){

	return /^([^\x00-\x20\x22\x28\x29\x2c\x2e\x3a-\x3c\x3e\x40\x5b-\x5d\x7f-\xff]+|\x22([^\x0d\x22\x5c\x80-\xff]|\x5c[\x00-\x7f])*\x22)(\x2e([^\x00-\x20\x22\x28\x29\x2c\x2e\x3a-\x3c\x3e\x40\x5b-\x5d\x7f-\xff]+|\x22([^\x0d\x22\x5c\x80-\xff]|\x5c[\x00-\x7f])*\x22))*\x40([^\x00-\x20\x22\x28\x29\x2c\x2e\x3a-\x3c\x3e\x40\x5b-\x5d\x7f-\xff]+|\x5b([^\x0d\x5b-\x5d\x80-\xff]|\x5c[\x00-\x7f])*\x5d)(\x2e([^\x00-\x20\x22\x28\x29\x2c\x2e\x3a-\x3c\x3e\x40\x5b-\x5d\x7f-\xff]+|\x5b([^\x0d\x5b-\x5d\x80-\xff]|\x5c[\x00-\x7f])*\x5d))*$/.test( email );	
}

hideErrorMessage();

$(function () {
    $('form').on('submit', function(e) {

        e.preventDefault();

        hideErrorMessage();

        var registerURL = 'http://localhost:8080/register';
        var data = $('form').serializeArray();
        var reduced = data.reduce((acc, {name, value}) => ({...acc, [name]: value}),{}); // form the object
        let textJSON = JSON.stringify(reduced)


        console.log(textJSON)
        if (reduced.password.length <7) {
            showErrorMessage("register-error","Пароль должен быть хотя бы 7 символов в длину")
        } else if (reduced.password!==reduced.repeat_password) {
            showErrorMessage("register-error","Пароли не совпадают")
        } else if (!isEmail(reduced.email)) {
            showErrorMessage("register-error", "Введите валидный адрес почты.")
        } else if (reduced.username.length < 4) {
            showErrorMessage("register-error","Имя пользователя должно быть хотя бы 4 символа в длину")
        } else {
            let regJSON = '{"username":"'+reduced.username+'","password":"'+reduced.password+'"}'
            sendRegisterRequest(registerURL, regJSON)
        }
        

    });
});

function sendRegisterRequest(url, json) {
    $.ajax({
        type: "POST",
        url: url,
        data: json,
        success: function(data)
        {
            console.log(data)
            sendAuthRequest(json)
        },
        error: function(jqXHR) {
            if (jqXHR.status === 400) {
                showErrorMessage("register-error","Пользователь с таким именем уже существует");
            } else {
                showErrorMessage("Internal server error.");
            }
        }
    });
}

function sendAuthRequest(json) {
    authRequestURL = "http://localhost:8080/auth"
    $.ajax({
        type: "POST",
        url: authRequestURL,
        data: json,
        success: function(data)
        {
            $.cookie('access_token',data.access_token, {path:'/'})
            window.location.replace("http://localhost:9080/feed.html")
        },
        error: function(jqXHR) {
            if (jqXHR.status === 400) {
            } else {
                showErrorMessage("register-error","Internal server error.");
            }
        }
    });
}