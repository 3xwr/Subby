function parseJwt (token) {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
};

if ($.cookie('access_token') === undefined) {
    window.location.replace("http://localhost:9080/login.html");
}

//alert($.cookie('access_token'))
var endpoint = 'http://localhost:8080/posts';
var user_data_endpoint = "http://localhost:8080/userdata"
let user_id = parseJwt($.cookie('access_token'))

user_json = JSON.parse('{"user_id":"'+user_id.sub+'"}')
user_json = JSON.stringify(user_json)
console.log(user_json)

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
       console.log(data)
       avatar=data.avatar_ref
       $('#right-navbar')
       .append(
        $("<h5>")
        .text(data.name)
        .attr("id","logged-in-user")
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
                       .text("Logout")
                   )
               )
           )
       )
       logoutClickListener();
    },
    error: function(jqXHR) {
        if (jqXHR.status === 403 || jqXHR.status === 401) {
            alert('403 or 401');
            window.location.replace("http://localhost:9080/login.html");
        }
    }
});


$.ajax({
    type: "GET",
    url: endpoint,
    headers: {
        "Authorization":"Bearer "+$.cookie('access_token')
    },
    success: function(data)
    {
        if (data===null) {
            $('#posts')
            .append(
                $("<div>")
                .addClass("font-monospace")
                .text("No posts in your feed. Subscribe to somebody!")
            )
        } else {
            console.log(data);
            $(document).ready(function(){
                for (let i = 0; i < data.length;i++) {
                    let img_src = data[i].image_ref
                    posted_at = data[i].posted_at
                    date = new Date(posted_at)
                    const formattedDate = date.toLocaleString("en-GB", {
                        day: "numeric",
                        month: "short",
                        year: "numeric",
                        hour: "numeric",
                        minute: "2-digit"
                      });
                    $('#posts')
                    .append(
                        $("<div>").addClass("card")
                        .append(
                            $("<div>").addClass("card-body")
                            .append(
                                $("<h5>").addClass("card-title d-flex justify-content-start").text(data[i].poster_username)
                            )
                            .append(
                                $("<h6>").text(formattedDate).addClass("card-subtitle mb-2 text-muted d-flex justify-content-start")
                            )
                            .append(
                                $("<p>").addClass("card-text").text(data[i].body)
                            )
                        )
                        .append(
                            $('<img>').attr("src", base_path+img_src).addClass("card-img-top")
                        )
                    )
                }
            });
        }        
    },
    error: function(jqXHR) {
        if (jqXHR.status === 403 || jqXHR.status === 401) {
            alert('403 or 401');
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