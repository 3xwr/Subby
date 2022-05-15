
const upload_url = "http://localhost:8080/upload"
const post_submit_url = "http://localhost:8080/post"

function submitForm(data) {
    console.log(data)
    if(data === undefined) {
        let textJSON = '{"body":"'+jsonData.post_body+'","membership_locked":false"}'
        $.ajax({
            type: "POST",
            url: post_submit_url,
            data: textJSON,
            headers: {
                "Authorization":"Bearer "+$.cookie('access_token')
            },
            success: function(data)
            {
                //
            },
            error: function(e) {
                //
            }
        });
    } else {
        let textJSON = '{"body":"'+jsonData.post_body+'","membership_locked":false,"image_ref":"'+data.file_address+'"}'
        $.ajax({
            type: "POST",
            url: post_submit_url,
            data: textJSON,
            headers: {
                "Authorization":"Bearer "+$.cookie('access_token')
            },
            success: function(data)
            {
                //
            },
            error: function(e) {
                $("#result").text(data);
                console.log("SUCCESS : ", data);
                $("#btnSubmit").prop("disabled", false);
            }
        });
    }
}

$(function () {
    $('form').on('submit', function(e) {

        e.preventDefault();

        var form = $('#upload-form')[0];

        var data = new FormData(form);

        $("#btnSubmit").prop("disabled", true);

        var formData = $('form').serializeArray();
        var reduced = formData.reduce((acc, {name, value}) => ({...acc, [name]: value}),{}); // form the object
        reduced = JSON.stringify(reduced)
        jsonData = JSON.parse(reduced)
        console.log(reduced)

        $.ajax({
            type: "POST",
            enctype: 'multipart/form-data',
            url: upload_url,
            headers: {
                "Authorization":"Bearer "+$.cookie('access_token')
            },
            data: data,
            processData: false,
            contentType: false,
            cache: false,
            timeout: 600000,
            success: function (data) {
                submitForm(data)
            },
            error: function (jqXHR) {
                if (jqXHR.status === 403 || jqXHR.status === 401) {
                    alert('403 or 401');
                    window.location.replace("http://localhost:9080/login.html");
                }
                if (jqXHR.status === 400) {
                    submitForm()
                }
            }
        });
    });
});

