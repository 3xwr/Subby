function hideErrorMessage(id) {
  if (id === undefined) {
    id = "password-error";
  }
  document.getElementById(id).style.display = "none";
}

function showErrorMessage(id, msg) {
  if (id === undefined) {
    id = "password-error";
  }
  document.getElementById(id).innerHTML = msg;
  document.getElementById(id).style.display = "inline-block";
}

function hideMessage(id) {
  document.getElementById(id).style.display = "none";
}

function showMessage(id) {
  document.getElementById(id).style.display = "block";
}

function parseJwt(token) {
  var base64Url = token.split(".")[1];
  var base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  var jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map(function (c) {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join("")
  );

  return JSON.parse(jsonPayload);
}

const upload_url = "http://localhost:8080/upload";
let img_address = "";
base_path = "http://localhost:9080/img/";

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

$(function () {
  hideErrorMessage();
  hideErrorMessage("avatar-error");
  hideMessage("avatar-ok");
  hideMessage("password-ok");
  $("form").on("submit", function (e) {
    e.preventDefault();
  });
});

async function getLoggedInUserID() {
  return new Promise((resolve) => {
    let checkPath = "http://localhost:8080/livecheck";
    if ($.cookie("access_token") === undefined) {
      resolve(false);
      window.location.replace("http://localhost:9080/login.html");
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
            console.log(logged_in_id);
            resolve(logged_in_id);
          } else {
            $.removeCookie("access_token", { path: "/" });
            window.location.replace("http://localhost:9080/login.html");
          }
        },
      });
    }
  });
}

async function getUserData(id) {
  console.log(id);
  return new Promise((resolve) => {
    let getDataPath = "http://localhost:8080/userprivatedata";
    getDataJSON = '{"user_id":"' + id + '","full_info":true}';
    $.ajax({
      type: "POST",
      url: getDataPath,
      data: getDataJSON,
      success: function (data) {
        console.log(data);
        resolve(data);
      },
      error: function (jqXHR) {},
    });
  });
}

function changePasswordListener(btnID) {
  changePasswordPath = "http://localhost:8080/changepassword";
  console.log("listening to " + btnID);
  $(btnID).on("click", function () {
    old_password = $("#oldPassword").val();
    console.log(old_password);
    new_password = $("#newPassword").val();
    console.log(new_password);
    repeat_password = $("#repeatPassword").val();
    console.log(repeat_password);

    if (
      old_password.length < 7 ||
      new_password.length < 7 ||
      repeat_password.length < 7
    ) {
      showErrorMessage(undefined, "Пароль должен быть длиннее 6 символов.");
    } else if (new_password !== repeat_password) {
      showErrorMessage(undefined, "Пароли не совпадают.");
    }

    changePasswordJSON =
      '{"old_password":"' +
      old_password +
      '","new_password":"' +
      new_password +
      '"}';
    console.log(changePasswordJSON);
    $.ajax({
      type: "POST",
      url: changePasswordPath,
      data: changePasswordJSON,
      headers: {
        Authorization: "Bearer " + $.cookie("access_token"),
      },
      success: function (data) {
        showMessage("password-ok");
        window.setTimeout(function () {
          document.location.reload();
        }, 1000);
      },
      error: function (jqXHR) {
        if (jqXHR.status === 403 || jqXHR.status === 401) {
          window.location.replace("http://localhost:9080/login.html");
        } else {
          showErrorMessage(undefined, "Неверный старый пароль.");
        }
      },
    });
  });
}

function changeAvatarListener(btnID) {
  changeAvatarPath = "http://localhost:8080/changeavatar";
  console.log("listening to " + btnID);
  $(btnID).on("click", function () {
    console.log(img_address);
    if (img_address=="") {
        showErrorMessage("avatar-error","Вы не загрузили новый аватар.")
    } else {
        hideErrorMessage("avatar-error")

        changeAvatarJSON = '{"avatar_ref":"'+img_address+'"}'
    $.ajax({
        type: "POST",
        url: changeAvatarPath,
        data: changeAvatarJSON,
        headers: {
          Authorization: "Bearer " + $.cookie("access_token"),
        },
        success: function (data) {
          showMessage("avatar-ok");
          window.setTimeout(function () {
            document.location.reload();
          }, 1000);
        },
        error: function (jqXHR) {
          if (jqXHR.status === 403 || jqXHR.status === 401) {
            window.location.replace("http://localhost:9080/login.html");
          } else {
            showErrorMessage(undefined, "Неверный старый пароль.");
          }
        },
      });
    }
  });
}

function changeDisabledFields(name, email) {
  $("#username-label").text(name);
  $("#email-label").text(email);
}

async function doStuff() {
  loggedInId = await getLoggedInUserID();
  userData = await getUserData(loggedInId);
  changeDisabledFields(userData.name, userData.email);
  changePasswordListener("#btn-change-password");
  changeAvatarListener("#btn-change-avatar")
}

doStuff();
