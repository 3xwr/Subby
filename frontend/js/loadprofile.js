function parseJwt (token) {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
};

var user_data_endpoint = "http://localhost:8080/userdata"
profile_base_path = "http://localhost:9080/user_profile.html?user="

let login_path = "/login.html"

if ($.cookie('access_token') === undefined && window.location.pathname!=login_path) {

} else {

    let user_id = parseJwt($.cookie('access_token'))
    
    user_json = JSON.parse('{"user_id":"'+user_id.sub+'"}')
    user_json = JSON.stringify(user_json)
    // console.log(user_json)
    
    let user_object
    
    base_path = "http://localhost:9080/img/"
    
    $.ajax({
        type: "POST",
        url: user_data_endpoint,
        data: user_json,
        headers: {
            "Content-type":"application/json"
        },
        success: function(data)
        {
        //    console.log(data)
           avatar=data.avatar_ref
           $('#right-navbar')
           .append(
            $("<a>")
            .append(
                $("<h5>")
                .text(data.name)
                .attr("id","logged-in-user")
               )
               .attr("href",profile_base_path+data.name)
            )
            
    
           .append(
               $("<div>").addClass("dropdown")
               .append(
                   $("<a>")
                   .addClass("dropdown-toggle d-flex align-items-center hidden-arrow")
                   .attr("id","navbarDropdownMenuAvatar")
                   .attr("role","button")
                   .attr("data-mdb-toggle","dropdown")
                   .attr("aria-expanded","false")
                   .append(
                       $('<img>')
                       .attr("src", base_path+avatar)
                       .addClass("rounded-circle")
                       .attr("height","30")
                       .attr("width", "30")
                       .attr("loading","lazy")
                   )
               )
               .append(
                   $("<ul>")
                   .addClass("dropdown-menu dropdown-menu-end")
                   .attr("aria-labelledby","navbarDropdownMenuAvatar")
                   .append(
                       $('<li>')
                       .append(
                           $('<a>')
                           .addClass("dropdown-item")
                           .attr("id","logout")
                           .text("Выход")
                       )
                   )
               )
           )
           logoutClickListener();
        },
        error: function(jqXHR) {
            if (jqXHR.status === 403 || jqXHR.status === 401) {
                window.location.replace("http://localhost:9080/login.html");
            }
        }
    });
    
    function logoutClickListener() {
        $("#logout").off();
    
        $("#logout").on("click", function(){
            $.removeCookie('access_token', {path:'/'});
            location.reload();
        })
    }
}


