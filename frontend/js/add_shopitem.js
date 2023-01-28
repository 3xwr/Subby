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
const add_item_url = "http://localhost:8080/additem"
let img_address = ''
base_path = "http://localhost:9080/img/"

let empty_description = true;
let empty_name = true;
let empty_price = true;

async function getLoggedInUserID() {
    return new Promise((resolve) => {
      let checkPath = "http://localhost:8080/livecheck";
      if ($.cookie("access_token") === undefined) {
        resolve(false);
      } else {
        $.ajax({
          type: "GET",
          url: checkPath,
          headers: {
            Authorization: "Bearer " + $.cookie("access_token"),
          },
          success: function (data) {
            //resolve(data)
          },
          error: function (jqXHR) {
            if (jqXHR.status === 404) {
              logged_in_id = parseJwt($.cookie("access_token")).sub;
              resolve(logged_in_id);
            } else {
              $.removeCookie("access_token", { path: "/" });
              document.location.reload();
            }
          },
        });
      }
    });
  }

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
           // 
            $("#preview").attr('src',base_path+img_address)
            $("#preview-label").removeAttr("hidden")
            showMessage("preview")
            if(!(empty_name&&empty_description&&empty_price)){
                $('.btn-primary').attr('disabled', false);
            }
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
    $('.form-outline #item-name').on('keyup', function() {
        $('.form-outline #item-name').each(function() {
          empty_name = !$.trim($("#item-name").val());
        });

        if((empty_description||empty_name||empty_price||img_address==''))
            $('#btnSubmitItem').attr('disabled', 'disabled');
        else
          $('#btnSubmitItem').attr('disabled', false);
        //  console.log(empty_name,empty_price,empty_rewards)
        console.log(empty_name,empty_price,empty_description, img_address)
      });
      $('.form-outline #item-price').on('keyup', function() {
        $('.form-outline #item-price').each(function() {
          empty_price = !$.trim($("#item-price").val());
        });

        if((empty_description||empty_name||empty_price||img_address==''))
            $('#btnSubmitItem').attr('disabled', 'disabled');
        else
          $('#btnSubmitItem').attr('disabled', false);
         // console.log(empty_name,empty_price,empty_rewards)
         console.log(empty_name,empty_price,empty_description, img_address)
      });
    $('.form-outline textarea').on('keyup', function() {
      $('.form-outline textarea').each(function() {
        empty_description = !$.trim($("#item-description-input").val());
      });
  
      if((empty_description||empty_name||empty_price||img_address=='')) {
        $('#btnSubmitItem').attr('disabled', 'disabled');
      }
      else {
        $('#btnSubmitItem').attr('disabled', false);
      }
        console.log(empty_name,empty_price,empty_description, img_address)
      });
  });


$(function () {
    $('form').on('submit', function(e) {

        e.preventDefault();

        getLoggedInUserID().then(
            function(currentUserID) {
                if (currentUserID===false) {
                    document.location.reload()
                } else {
                    var form = $('#upload-form')[0];

                    var data = new FormData(form);
            
                    $("#btnSubmitItem").prop("disabled", true);
            
                    var formData = $('form').serializeArray();
                    console.log(formData)
                    var reduced = formData.reduce((acc, {name, value}) => ({...acc, [name]: value}),{}); // form the object
                    reduced.price = parseInt(reduced.price, 10)
                    let textJSON = JSON.stringify(reduced)
                    console.log(reduced.name)
            
                    textJSON = '{"owner_id":"'+currentUserID+'","name":"'+reduced.name+'","price":'+reduced.price+',"description":"'+reduced.description+'","image_ref":"'+img_address+'"}'
                    console.log(textJSON)
                    $.ajax({
                        type: "POST",
                        url: add_item_url,
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
                            $("#btnSubmitItem").prop("disabled", false);
                        }
                    });
                }
            });
    });
});

