function hideMessage(id) {
    document.getElementById(id).style.display ="none";
}

function showMessage(id) {
    document.getElementById(id).style.display = "block";
}

token = $.cookie('access_token')
if (token===undefined) {
    window.location.replace("http://localhost:9080/login.html")
}

const upload_url = "http://localhost:8080/upload"
const add_tier_url = "http://localhost:8080/addtier"
let img_address = ''
base_path = "http://localhost:9080/img/"

let empty_rewards = true;
let empty_name = true;
let empty_price = true;

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
           // $('.btn-primary').attr('disabled', false);
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
    $('.form-outline #tier-name').on('keyup', function() {
        $('.form-outline #tier-name').each(function() {
          empty_name = !$.trim($("#tier-name").val());
        });

        if((empty_rewards||empty_name||empty_price))
            $('#btnSubmitTier').attr('disabled', 'disabled');
        else
          $('#btnSubmitTier').attr('disabled', false);
        //  console.log(empty_name,empty_price,empty_rewards)
      });
      $('.form-outline #tier-price').on('keyup', function() {
        $('.form-outline #tier-price').each(function() {
          empty_price = !$.trim($("#tier-price").val());
        });

        if((empty_rewards||empty_name||empty_price))
            $('#btnSubmitTier').attr('disabled', 'disabled');
        else
          $('#btnSubmitTier').attr('disabled', false);
         // console.log(empty_name,empty_price,empty_rewards)
      });
    $('.form-outline textarea').on('keyup', function() {
      $('.form-outline textarea').each(function() {
        empty_rewards = !$.trim($("#tier-rewards-input").val());
      });
  
      if((empty_rewards||empty_name||empty_price)) {
        $('#btnSubmitTier').attr('disabled', 'disabled');
      }
      else {
        $('#btnSubmitTier').attr('disabled', false);
      }
        //console.log(empty_name,empty_price,empty_rewards)
      });
  });


$(function () {
    $('form').on('submit', function(e) {

        e.preventDefault();

        var form = $('#upload-form')[0];

        var data = new FormData(form);

        $("#btnSubmitTier").prop("disabled", true);

        var formData = $('form').serializeArray();
        console.log(formData)
        var reduced = formData.reduce((acc, {name, value}) => ({...acc, [name]: value}),{}); // form the object
        reduced.price = parseInt(reduced.price, 10)
        textJSON = JSON.stringify(reduced)
        console.log(reduced.name)

        if (img_address==='') {
            $.ajax({
                type: "POST",
                url: add_tier_url,
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
            let textJSON = '{"name":"'+reduced.name+'","price":'+reduced.price+',"rewards":"'+reduced.rewards+'","image_ref":"'+img_address+'"}'
            console.log(textJSON)
            $.ajax({
                type: "POST",
                url: add_tier_url,
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

