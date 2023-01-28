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

shop_base_path = "http://localhost:9080/user_shop.html?user=";

let username = findGetParameter("user");
if (username === null || username === "") {
  window.location.replace("http://localhost:9080/feed.html");
}

userIDByNamePath = "http://localhost:8080/useridbyname";
userIDRequestJSON = JSON.parse('{"username":"' + username + '"}');
userIDRequestJSON = JSON.stringify(userIDRequestJSON);

let token = $.cookie("access_token");
var userPageID;
var membershipID;
let subbed = true;

console.log(username);

async function getUserIDByUsername(username) {
  return new Promise((resolve, reject) => {
    userIDByNamePath = "http://localhost:8080/useridbyname";
    userIDRequestJSON = JSON.parse('{"username":"' + username + '"}');
    userIDRequestJSON = JSON.stringify(userIDRequestJSON);
    $.ajax({
      type: "POST",
      url: userIDByNamePath,
      data: userIDRequestJSON,
      headers: {
        "Content-type": "application/json",
      },
      success: function (data) {
        resolve(data.user_id);
      },
    });
  });
}

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
            console.log(logged_in_id);
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

async function getUserShopByOwnerID(id) {
  let getShopPath = "http://localhost:8080/shop";
  let shopRequestJSON = '{"owner_id":"' + id + '"}';
  return new Promise((resolve) => {
    $.ajax({
      type: "POST",
      url: getShopPath,
      data: shopRequestJSON,
      success: function (data) {
        resolve(data);
      },
      error: function (jqXHR) {},
    });
  });
}

function buildDivs(shopData) {
  $("#shop-content")
  .append($("<div>").addClass("row"))
  .before(
    $("<h2>")
      .text("Магазин " + username)
      .attr("id", "shop-title")
  );
  if (shopData === null ){
    $("#shop-title")
    .append(
      $("<h5>")
      .text("У пользователя нет товаров.")
      .addClass("text-muted")
      .attr("id","no-items")
    )
  } else {
    for (let index = 0; index < shopData.length; index++) {
      const element = shopData[index];
      $(".row").append(buildItemCard(element,index));
    }
  }


}

function buildItemCard(item,id) {
  baseImgPath = "http://localhost:9080/img/";
  let item_div = $("<div>")
    .addClass("col-md-3")
    .attr("id","column")
    .append(
      $("<div>")
      .addClass("card")
      .append(
        $("<img>")
          .attr("src", baseImgPath + item.image_ref)
          .addClass("card-img-top")
      )
      .append(
        $("<div>")
          .addClass("card-body")
          .append($("<h5>").addClass("card-title").text(item.name))
          .append($("<p>").addClass("card-text").text(item.description)
      )
      )
      .append(
        $("<div>")
        .addClass("card-footer")
        .attr("id","item"+id)
        .append(
          $("<a>").addClass("btn btn-primary").text(item.price+" ₽").attr("id","buyBtn")
        )
      )
    )
    getLoggedInUserID().then(
      function(value) {
        if(value!==false){
          if (item.owner_id == value) {
            crossPath = "http://localhost:9080/img/small-cross.png"
            $("#item" + id).after(
                $("<a>")
                .attr("id", "deletable-item"+id)
                .addClass("deletable")
                .append(
                    $("<img>").attr("src", crossPath).attr("id","delete-img")
                )
              )
              deleteItemListener("#deletable-item"+id,item.id)
        }
        }
      }
    )
    
    
  return item_div;
}

function deleteItemListener(btnID, itemID) {
  deleteItemPath = "http://localhost:8080/deleteitem";
  console.log("listening to "+btnID+", id:"+itemID)
  $(btnID).on("click", function () {

    userUUID = parseJwt($.cookie("access_token"));
    deleteItemJSON ='{"item_id":"'+itemID+'"}';
    $.ajax({
      type: "POST",
      url: deleteItemPath,
      data: deleteItemJSON,
      headers: {
        Authorization: "Bearer " + $.cookie("access_token"),
      },
      success: function (data) {
        document.location.reload();
      },
      error: function (jqXHR) {
        if (jqXHR.status === 403 || jqXHR.status === 401) {
          window.location.replace("http://localhost:9080/login.html");
        }
      },
    });
  });
}

async function doStuff() {
  shopUserID = await getUserIDByUsername(username);
  loggedInUserID = await getLoggedInUserID();
  shopData = await getUserShopByOwnerID(shopUserID);
  console.log("SHOP USER ID", shopUserID);
  console.log("LOGGED IN USER ID", loggedInUserID);
  console.log("shop data", shopData);
  buildDivs(shopData);
}

doStuff();
