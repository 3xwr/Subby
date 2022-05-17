function hideMessage(id) {
    document.getElementById(id).style.display ="none";
}

function showMessage(id) {
    document.getElementById(id).style.display = "block";
}

const upload_url = "http://localhost:8080/upload"
const post_submit_url = "http://localhost:8080/post"
let img_address = ''
base_path = "http://localhost:9080/img/"

let empty = false;

$("input[name='img']").change(function() {

    var form = $('#upload-form')[0];

    var data = new FormData(form);

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
            img_address = data.file_address
            $('.btn-primary').attr('disabled', false);
            $("#preview").attr('src',base_path+img_address)
            $("#preview-label").removeAttr("hidden")
            showMessage("preview")
        },
        error: function (jqXHR) {
            if (jqXHR.status === 403 || jqXHR.status === 401) {
                window.location.replace("http://localhost:9080/login.html");
            }
            if (jqXHR.status === 400) {
                //
            }
        }
    });
})

$(document).ready(function() {
    $('.form-outline textarea').on('keyup', function() {
  
      $('.form-outline textarea').each(function() {
        empty = !$.trim($("#post_body_input").val());
      });
  
      if (empty && img_address==='')
        $('.btn-primary').attr('disabled', 'disabled');
      else
        $('.btn-primary').attr('disabled', false);
    });
  });


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

        if (img_address==='') {
            let textJSON = '{"body":"'+jsonData.post_body+'","membership_locked":false}'
            $.ajax({
                type: "POST",
                url: post_submit_url,
                data: textJSON,
                headers: {
                    "Authorization":"Bearer "+$.cookie('access_token')
                },
                success: function(data)
                {
                    showMessage("submit-success")
                    window.setTimeout(function(){
                        window.location.replace("http://localhost:9080/feed.html");
                    },1000);
                },
                error: function(e) {
                    $("#result").text(data);
                    console.log("SUCCESS : ", data);
                    $("#btnSubmit").prop("disabled", false);
                }
            });
        } else {
            let textJSON = '{"body":"'+jsonData.post_body+'","membership_locked":false,"image_ref":"'+img_address+'"}'
            $.ajax({
                type: "POST",
                url: post_submit_url,
                data: textJSON,
                headers: {
                    "Authorization":"Bearer "+$.cookie('access_token')
                },
                success: function(data)
                {
                    showMessage("submit-success")
                    window.setTimeout(function(){
                        window.location.replace("http://localhost:9080/feed.html");
                    },1000);
                },
                error: function(e) {
                    $("#result").text(data);
                    console.log("SUCCESS : ", data);
                    $("#btnSubmit").prop("disabled", false);
                }
            });
        }
    });
});

