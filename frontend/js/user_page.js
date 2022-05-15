var endpoint = 'http://localhost:8080/posts';

base_path = "http://localhost:9080/img/"

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
                    if(img_src===undefined) {
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
                        )
                    } else {
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
