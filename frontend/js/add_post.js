function hideMessage(id) {
  document.getElementById(id).style.display = "none";
}

function showMessage(id) {
  document.getElementById(id).style.display = "block";
}

token = $.cookie("access_token");
if (token === undefined) {
  window.location.replace("http://localhost:9080/login.html");
}

const upload_url = "http://localhost:8080/upload";
const post_submit_url = "http://localhost:8080/post";
let img_address = "";
base_path = "http://localhost:9080/img/";

let empty = false;

getUserTiers();

$("input[name='img']").change(function () {
  var form = $("#upload-form")[0];

  var data = new FormData(form);

  $.ajax({
    type: "POST",
    enctype: "multipart/form-data",
    url: upload_url,
    headers: {
      Authorization: "Bearer " + $.cookie("access_token"),
    },
    data: data,
    processData: false,
    contentType: false,
    cache: false,
    timeout: 600000,
    success: function (data) {
      img_address = data.file_address;
      $("#btnSubmitPost").attr("disabled", false);
      $("#preview").attr("src", base_path + img_address);
      $("#preview-label").removeAttr("hidden");
      showMessage("preview");
    },
    error: function (jqXHR) {
      if (jqXHR.status === 403 || jqXHR.status === 401) {
        window.location.replace("http://localhost:9080/login.html");
      }
      if (jqXHR.status === 400) {
        //
      }
    },
  });
});

$(document).ready(function () {
  $(".form-outline textarea").on("keyup", function () {
    $(".form-outline textarea").each(function () {
      empty = !$.trim($("#post_body_input").val());
    });

    if (empty && img_address === "")
      $("#btnSubmitPost").attr("disabled", "disabled");
    else $("#btnSubmitPost").attr("disabled", false);
  });
});

$(function () {
  $("form").on("submit", function (e) {
    e.preventDefault();

    var form = $("#upload-form")[0];

    var data = new FormData(form);

    $("#btnSubmitPost").prop("disabled", true);

    var formData = $("form").serializeArray();
    var reduced = formData.reduce(
      (acc, { name, value }) => ({ ...acc, [name]: value }),
      {}
    ); // form the object
    reduced = JSON.stringify(reduced);
    jsonData = JSON.parse(reduced);

    let membership_locked = false;
    let values;

    if ($("#paid-post-check").is(":checked")) {
      const selected = document.querySelectorAll("#tier-select option:checked");
      arr = Array.from(selected).map((el) => el.value);
      membership_locked = true;
      membership_tier = arr;
      values = JSON.stringify(arr);
      console.log(arr);
      console.log(arr.length)
      //temporary fix, should instead disable the send button if checkbox is clicked but no tiers are selected
      if (arr.length == 0) {
        membership_locked = false;
        $("#paid-post-check").prop("checked", false);
        $("#tier-select").next().hide();
      }
    }

    if (img_address === "") {
      var textJSON;
      if (membership_locked) {
        textJSON =
          '{"body":"' +
          jsonData.post_body +
          '","membership_locked":' +
          membership_locked +
          ',"membership_tiers":' +
          values +
          "}";
      } else {
        textJSON =
          '{"body":"' +
          jsonData.post_body +
          '","membership_locked":' +
          membership_locked +
          "}";
      }
      $.ajax({
        type: "POST",
        url: post_submit_url,
        data: textJSON,
        headers: {
          Authorization: "Bearer " + $.cookie("access_token"),
        },
        success: function (data) {
          showMessage("submit-success");
          window.setTimeout(function () {
            window.location.replace("http://localhost:9080/feed.html");
          }, 1000);
        },
        error: function (e) {
          $("#result").text(data);
          console.log("SUCCESS : ", data);
          $("#btnSubmitPost").prop("disabled", false);
        },
      });
    } else {
      var textJSON;
      if (membership_locked) {
        textJSON =
          '{"body":"' +
          jsonData.post_body +
          '","membership_locked":' +
          membership_locked +
          ',"membership_tiers":' +
          values +
          ',"image_ref":"' +
          img_address +
          '"}';
      } else {
        textJSON =
          '{"body":"' +
          jsonData.post_body +
          '","membership_locked":false,"image_ref":"' +
          img_address +
          '"}';
      }
      $.ajax({
        type: "POST",
        url: post_submit_url,
        data: textJSON,
        headers: {
          Authorization: "Bearer " + $.cookie("access_token"),
        },
        success: function (data) {
          showMessage("submit-success");
          window.setTimeout(function () {
            window.location.replace("http://localhost:9080/feed.html");
          }, 1000);
        },
        error: function (e) {
          $("#result").text(data);
          console.log("SUCCESS : ", data);
          $("#btnSubmitPost").prop("disabled", false);
        },
      });
    }
  });
});

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

//get user membership tiers
async function getMembershipIDByOwnerID(id) {
  return new Promise((resolve, reject) => {
    let membershipOwnerPath = "http://localhost:8080/membershipowner";
    let ownerIDJSON = '{"owner_id":"' + id + '"}';
    $.ajax({
      type: "POST",
      url: membershipOwnerPath,
      data: ownerIDJSON,
      success: function (data) {
        resolve(data.membership_id);
      },
      error: function () {
        resolve(false);
      },
    });
  });
}

async function getMembershipDataByID(id) {
  return new Promise((resolve, reject) => {
    let membershipPath = "http://localhost:8080/membership";
    let membershipIDJSON = '{"membership_id":"' + id + '"}';
    $.ajax({
      type: "POST",
      url: membershipPath,
      data: membershipIDJSON,
      success: function (data) {
        resolve(data);
      },
      error: function () {
        resolve(false);
      },
    });
  }).catch();
}

async function getUserTiers() {
  owner_id = await getLoggedInUserID();
  membership_id = await getMembershipIDByOwnerID(owner_id);
  if (membership_id === false) {
    hasData = membership_id;
    buildTierDivs(hasData, hasData);
  } else {
    membership_data = await getMembershipDataByID(membership_id);
    buildTierDivs(membership_data, true);
    console.log(membership_data);
  }
}

function buildTierDivs(membership_data, hasData) {
  if (!hasData) {
    console.log(hasData);
    $("#tier-select").hide();
    $(".form-check-input").prop("disabled", true);
  } else {
    if (membership_data.tiers !== null) {
      $('#tier-select').multiSelect({
        selectableHeader: "<div class='custom-header'>Доступные уровни</div>",
        selectionHeader: "<div class='custom-header'>Выбранные уровни</div>"
      });
      $("#tier-select").multiSelect({ keepOrder: true });
      for (let i = 0; i < membership_data.tiers.length; i++) {
        tier = membership_data.tiers[i];
        $("#tier-select").multiSelect("addOption", {
          value: tier.id,
          text: tier.name,
        });
      }
      $("#tier-select").next().hide();
      $(".form-check-input").on("click", function () {
        var checkbox = $(this);
  
        if (checkbox.is(":checked")) {
          $("#tier-select").next().show(300);
        } else {
          $("#tier-select").next().hide(200);
        }
      });
    } else {
      $("#tier-select").hide();
      $(".form-check-input").prop("disabled", true);
    }
 
  }
}
