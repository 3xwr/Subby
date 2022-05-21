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
      $("<p>")
        .text("Магазин " + username)
        .attr("id", "shop-title")
    );
  for (let index = 0; index < shopData.length; index++) {
    const element = shopData[index];
    $(".row").append(buildItemCard(element));
  }
}

function buildItemCard(item) {
  baseImgPath = "http://localhost:9080/img/";
  let item_div = $("<div>")
    .addClass("col-md-3")
    .attr("id","column")
    // .append(
    //   $("<div>")
    //   .text("COLUMN TEXT")
    // )
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
          .append(
            $("<h6>").addClass("card-subtitle mb-2 text-muted").text(item.price)
          )
          .append($("<p>").addClass("card-text").text(item.description))
      )
    )
    
  return item_div;
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
