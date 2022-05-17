function findGetParameter(parameterName) {
    var result = null,
        tmp = [];
    location.search
        .substr(1)
        .split("&")
        .forEach(function (item) {
          tmp = item.split("=");
          if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
        });
    return result;
}



profile_base_path = "http://localhost:9080/user_profile.html?user="

let username = findGetParameter("user")
if (username===null || username==="") {
    window.location.replace("http://localhost:9080/feed.html");
}


userIDByNamePath = "http://localhost:8080/useridbyname"
userIDRequestJSON = JSON.parse('{"username":"'+username+'"}')
userIDRequestJSON = JSON.stringify(userIDRequestJSON)

let token = $.cookie('access_token')
var userPageID
let subbed = true


function changeSubButton() {
    if (token !== undefined) {
        checkUserSubPath = "http://localhost:8080/checksub"
        userID = parseJwt(token)
        
        $.ajax({
            type: "POST",
            url: userIDByNamePath,
            data: userIDRequestJSON,
            headers: {
                "Content-type":"application/json"
            },
            success: function(data)
            {
                checkJSON = '{"subscriber":"'+userID.sub+'","subscribed":"'+data.user_id+'"}'
                userPageID=data.user_id
                $.ajax({
                    type: "POST",
                    url: checkUserSubPath,
                    data: checkJSON,
                    headers: {
                        "Authorization":"Bearer "+$.cookie('access_token')
                    },
                    success: function(subdata)
                    {
                        let userProfileRequestBody = '{"user_id":"'+data.user_id+'","full_info":true}'
                        $.ajax({
                            type: "POST",
                            url: user_data_endpoint,
                            data: userProfileRequestBody,
                            headers: {
                                "Content-type":"application/json"
                            },
                            success: function(subcount)
                            {
                                $("#sub_counter")
                                .text("Подписчики: " + subcount.subscriber_count)
                            },
                            error: function(jqXHR) {
                            }
                        });
                        if (subdata.subscribed == false) {
                            $("#sub_button")
                            .text("Подписаться")
                            subbed = false
                        } else {
                            $("#sub_button")
                            .text("Отписаться")
                            subbed = true
                        }
                    },
                    error: function(jqXHR) {
                    }
                });
            },
            error: function(jqXHR) {
            }
        });
    
    }
}


//get user profile
$.ajax({
    type: "POST",
    url: userIDByNamePath,
    data: userIDRequestJSON,
    headers: {
        "Content-type":"application/json"
    },
    success: function(data)
    {
        let userProfileRequestBody = '{"user_id":"'+data.user_id+'","full_info":true}'
        userPageID=data.user_id

        $.ajax({
            type: "POST",
            url: user_data_endpoint, //here is an error
            data: userProfileRequestBody,
            headers: {
                "Content-type":"application/json"
            },
            success: function(data)
            {
                base_path = "http://localhost:9080/img/"
                if (token === undefined || parseJwt(token).sub===userPageID) {
                    $("#userInfo")
                    .append(
                        $("<div>")
                        .addClass("p-5 text-center bg-light")
                        .append(
                            $("<img>")
                            .attr("src",base_path+data.avatar_ref)
                            .addClass("rounded-circle")
                            .attr("height","200")
                            .attr("width", "200")
                            .attr("loading","lazy")
                        )
                        .append(
                            $("<h1>")
                            .addClass("mb-3")
                            .text(data.name)
                        )
                        .append(
                            $("<h4>")
                            .addClass("mb-3")
                            .text("Подписчики: " + data.subscriber_count)
                            .attr("id","sub_counter")
                        )
                    )
                } else {
                    $("#userInfo")
                    .append(
                        $("<div>")
                        .addClass("p-5 text-center bg-light")
                        .append(
                            $("<img>")
                            .attr("src",base_path+data.avatar_ref)
                            .addClass("rounded-circle")
                            .attr("height","200")
                            .attr("width", "200")
                            .attr("loading","lazy")
                        )
                        .append(
                            $("<h1>")
                            .addClass("mb-3")
                            .text(data.name)
                        )
                        .append(
                            $("<h4>")
                            .addClass("mb-3")
                            .text("Подписчики: " + data.subscriber_count)
                            .attr("id","sub_counter")
                        )
                        .append(
                            $("<a>")
                            .addClass("btn btn-primary")
                            .attr("id","sub_button")
                        )
                    )
                    changeSubButton();
                    subscribeButtonClickListener()
                }
                var userPostsEndpoint = 'http://localhost:8080/userposts';
                let userPostsRequest = '{"poster_id":"'+userPageID+'"}'
                $.ajax({
                    type: "POST",
                    url: userPostsEndpoint,
                    data:userPostsRequest,
                    success: function(data)
                    {    
                        console.log(data)
                        for (let i = 0; i < data.length;i++) {
                            let img_src = data[i].image_ref
                            posted_at = data[i].posted_at
                            date = new Date(posted_at)
                            const formattedDate = date.toLocaleString("ru", {
                                day: "numeric",
                                month: "short",
                                year: "numeric",
                                hour: "numeric",
                                minute: "2-digit"
                              });
                            if(img_src===undefined) {
                                $('#posts')
                                .append(
                                    $("<div>").addClass("card")
                                    .append(
                                        $("<div>").addClass("card-body")
                                        .append(
                                            $("<a>")
                                            .addClass("d-flex justify-content-start")
                                            .attr("id","poster-data")
                                            .attr("href",profile_base_path+data[i].poster_username)
                                            .append(
                                                $("<img>")
                                                .attr("src",base_path+data[i].poster_avatar)
                                                .addClass("rounded-circle")
                                                .attr("height","30")
                                                .attr("width", "30")
                                                .attr("loading","lazy")
                                            )
                                            .append(
                                                $("<h5>")
                                                .text(data[i].poster_username)
                                                .attr("id","post-name")
                                            )
                                        )
                                        .append(
                                            $("<h6>").text(formattedDate).addClass("card-subtitle mb-2 text-muted d-flex justify-content-start")
                                        )
                                        .append(
                                            $("<p>").addClass("card-text").text(data[i].body)
                                        )
                                    )
                                )
                            } else {
                                $('#posts')
                                .append(
                                    $("<div>").addClass("card")
                                    .append(
                                        $("<div>").addClass("card-body")
                                        .append(
                                            $("<a>")
                                            .addClass("d-flex justify-content-start")
                                            .attr("id","poster-data")
                                            .attr("href",profile_base_path+data[i].poster_username)
                                            .append(
                                                $("<img>")
                                                .attr("src",base_path+data[i].poster_avatar)
                                                .addClass("rounded-circle")
                                                .attr("height","30")
                                                .attr("width", "30")
                                                .attr("loading","lazy")
                                            )
                                            .append(
                                                $("<h5>")
                                                .text(data[i].poster_username)
                                                .attr("id","post-name")
                                            )
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
                           
                        }
                    },
                    error: function(jqXHR) {
                        console.log(jqXHR)
                        if (jqXHR.status === 403 || jqXHR.status === 401) {
                            alert('403 or 401');
                            window.location.replace("http://localhost:9080/login.html");
                        }
                    }
                });
            },
            error: function(jqXHR) {
            }
        });
    },
    error: function(jqXHR) {
    }
});

function subscribeButtonClickListener() {
    $("#sub_button").off();

    subscribePath = "http://localhost:8080/subscribe"
    unsubscribePath = "http://localhost:8080/unsubscribe"

    $("#sub_button").on("click", function(){
        unsubJSON = '{"id":"'+userPageID+'"}'
        if (subbed) {
            $.ajax({
                type: "POST",
                url: unsubscribePath,
                data: unsubJSON,
                headers: {
                    "Authorization":"Bearer "+$.cookie('access_token')
                },
                success: function(data)
                {
                    changeSubButton();
                },
                error: function(jqXHR) {
                    if (jqXHR.status === 403 || jqXHR.status === 401) {
                        alert('403 or 401');
                        window.location.replace("http://localhost:9080/login.html");
                    }
                }
            });
        } else {
            $.ajax({
                type: "POST",
                url: subscribePath,
                data: unsubJSON,
                headers: {
                    "Authorization":"Bearer "+$.cookie('access_token')
                },
                success: function(data)
                {
                    changeSubButton();
                },
                error: function(jqXHR) {
                    if (jqXHR.status === 403 || jqXHR.status === 401) {
                        alert('403 or 401');
                        window.location.replace("http://localhost:9080/login.html");
                    }
                }
            });
        }

    })
}