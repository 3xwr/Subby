function hideErrorMessage() {
    document.getElementById("login-error").style.display ="none";
}

function showErrorMessage(msg) {
    document.getElementById("login-error").innerHTML = msg;
    document.getElementById("login-error").style.display = "inline-block";
}

hideErrorMessage();

$(function () {
    $('form').on('submit', function(e) {

        e.preventDefault();

        var endpoint = 'http://localhost:8080/auth';
        var data = $('form').serializeArray();
        var reduced = data.reduce((acc, {name, value}) => ({...acc, [name]: value}),{}); // form the object
        reduced = JSON.stringify(reduced)

        //console.log(reduced)
        
        $.ajax({
            type: "POST",
            url: endpoint,
            data: reduced,
            success: function(data)
            {
                //console.log(data.access_token);
                hideErrorMessage();
                //document.cookie='access_token='+data.access_token+';domain=;path=/'
                $.cookie('access_token',data.access_token, {path:'/'})
                window.location.replace("http://localhost:9080/feed.html")
            },
            error: function(jqXHR) {
                if (jqXHR.status === 403) {
                    showErrorMessage("Wrong username or password.");
                } else {
                    showErrorMessage("Internal server error.");
                }
            }
        });
    });
});