endpoint="http:localhost:8080/testlive"
if ($.cookie('access_token') == null) {
    $("#right-navbar")
    .append(
        $("<ul>")
        .addClass("navbar-nav me-auto mb-2 mb-lg-0")
        .append(
            $("<li>")
            .addClass("nav-item")
            .append(
                $("<a>")
                .addClass("nav-link")
                .attr("href","http://localhost:9080/login.html")
                .text("Войти")
            )
        )
    )
    $('.only-logged-in')
    .remove(
    )
} else {
    $.ajax({
        type: "GET",
        url: endpoint,
        success: function(data)
        {
        },
        error: function(jqXHR) {
            if (jqXHR.status === 404 || jqXHR.status === 401) {
            } else {
                $.removeCookie('access_token', {path:'/'});
                document.location.reload();
            }

        }
    });
}