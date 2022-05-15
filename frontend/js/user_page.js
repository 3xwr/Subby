//alert($.cookie('access_token'))
var endpoint = 'http://localhost:8080/posts';



$.ajax({
    type: "GET",
    url: endpoint,
    headers: {
        "Authorization":"Bearer "+$.cookie('access_token')
    },
    success: function(data)
    {
        console.log(data);
        $(document).ready(function(){
            for (let i = 0; i < data.length;i++) {
                let img_src = data[i].image_ref
                base_path = "http://localhost:9080/img/"
                $('#posts')
                .append(
                    $("<div>").addClass("card")
                    .append(
                        $("<div>").addClass("card-body")
                        .append(
                            $("<h5>").addClass("card-title d-flex justify-content-start").text(data[i].poster_username).attr("id","poster-username")
                        )
                        .append(
                            $("<h6>").text(data[i].posted_at).addClass("card-subtitle mb-2 text-muted d-flex justify-content-start")
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
    },
    error: function(jqXHR) {
        if (jqXHR.status === 403 || jqXHR.status === 401) {
            alert('403 or 401');
            window.location.replace("http://localhost:9080/login.html");
        }
    }
});