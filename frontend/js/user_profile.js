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

profile_base_path = "http://localhost:9080/user_profile.html?user=";

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

function changeSubButton() {
  if (token !== undefined) {
    checkUserSubPath = "http://localhost:8080/checksub";
    userID = parseJwt(token);

    $.ajax({
      type: "POST",
      url: userIDByNamePath,
      data: userIDRequestJSON,
      headers: {
        "Content-type": "application/json",
      },
      success: function (data) {
        checkJSON =
          '{"subscriber":"' +
          userID.sub +
          '","subscribed":"' +
          data.user_id +
          '"}';
        userPageID = data.user_id;
        $.ajax({
          type: "POST",
          url: checkUserSubPath,
          data: checkJSON,
          headers: {
            Authorization: "Bearer " + $.cookie("access_token"),
          },
          success: function (subdata) {
            let userProfileRequestBody =
              '{"user_id":"' + data.user_id + '","full_info":true}';
            $.ajax({
              type: "POST",
              url: user_data_endpoint,
              data: userProfileRequestBody,
              headers: {
                "Content-type": "application/json",
              },
              success: function (subcount) {
                $("#sub_counter").text(
                  "Подписчики: " + subcount.subscriber_count
                );
              },
              error: function (jqXHR) {},
            });
            if (subdata.subscribed == false) {
              $("#sub_button").text("Подписаться");
              subbed = false;
            } else {
              $("#sub_button").text("Отписаться");
              subbed = true;
            }
          },
          error: function (jqXHR) {},
        });
      },
      error: function (jqXHR) {},
    });
  }
}

//get user profile
$.ajax({
  type: "POST",
  url: userIDByNamePath,
  data: userIDRequestJSON,
  headers: {
    "Content-type": "application/json",
  },
  success: function (data) {
    let userProfileRequestBody =
      '{"user_id":"' + data.user_id + '","full_info":true}';
    userPageID = data.user_id;

    $.ajax({
      type: "POST",
      url: user_data_endpoint, //here is an error
      data: userProfileRequestBody,
      headers: {
        "Content-type": "application/json",
      },
      success: function (data) {
        base_path = "http://localhost:9080/img/";
        if (token === undefined || parseJwt(token).sub === userPageID) {
          $("#userInfo").append(
            $("<div>")
              .addClass("p-5 text-center bg-light")
              .append(
                $("<img>")
                  .attr("src", base_path + data.avatar_ref)
                  .addClass("rounded-circle")
                  .attr("height", "200")
                  .attr("width", "200")
                  .attr("loading", "lazy")
              )
              .append($("<h1>").addClass("mb-3").text(data.name))
              .append(
                $("<h4>")
                  .addClass("mb-3")
                  .text("Подписчики: " + data.subscriber_count)
                  .attr("id", "sub_counter")
              )
          );
        } else {
          $("#userInfo").append(
            $("<div>")
              .addClass("p-5 text-center bg-light")
              .append(
                $("<img>")
                  .attr("src", base_path + data.avatar_ref)
                  .addClass("rounded-circle")
                  .attr("height", "200")
                  .attr("width", "200")
                  .attr("loading", "lazy")
              )
              .append($("<h1>").addClass("mb-3").text(data.name))
              .append(
                $("<h4>")
                  .addClass("mb-3")
                  .text("Подписчики: " + data.subscriber_count)
                  .attr("id", "sub_counter")
              )
              .append(
                $("<a>").addClass("btn btn-primary").attr("id", "sub_button")
              )
          );
          changeSubButton();
          subscribeButtonClickListener();
        }
        var userPostsEndpoint = "http://localhost:8080/userposts";
        let userPostsRequest = '{"poster_id":"' + userPageID + '"}';
        getLoggedInUserID().then(
          function (value) {
            if (value === false) {
              $.ajax({
                type: "POST",
                url: userPostsEndpoint,
                data: userPostsRequest,
                success: function (data) {
                  console.log(data);
                  if (data == null) {
                    $("#posts").append(
                      $("<p>").text("У пользователя еще нет записей")
                    );
                  } else {
                    for (let i = 0; i < data.length; i++) {
                      let post_id = "post_" + i;
                      let img_src = data[i].image_ref;
                      posted_at = data[i].posted_at;
                      date = new Date(posted_at);
                      const formattedDate = date.toLocaleString("ru", {
                        day: "numeric",
                        month: "short",
                        year: "numeric",
                        hour: "numeric",
                        minute: "2-digit",
                      });
                      if (img_src === undefined) {
                        $("#posts").append(
                          $("<div>")
                            .addClass("card")
                            .append(
                              $("<div>")
                                .addClass("card-body")
                                .append(
                                  $("<a>")
                                    .addClass("d-flex justify-content-start")
                                    .attr("id", "poster-data")
                                    .attr(
                                      "href",
                                      profile_base_path +
                                        data[i].poster_username
                                    )
                                    .append(
                                      $("<img>")
                                        .attr(
                                          "src",
                                          base_path + data[i].poster_avatar
                                        )
                                        .addClass("rounded-circle")
                                        .attr("height", "30")
                                        .attr("width", "30")
                                        .attr("loading", "lazy")
                                    )
                                    .append(
                                      $("<h5>")
                                        .text(data[i].poster_username)
                                        .attr("id", "post-name")
                                    )
                                    .attr("id", post_id)
                                )
                                .append(
                                  $("<h6>")
                                    .text(formattedDate)
                                    .addClass(
                                      "card-subtitle mb-2 text-muted d-flex justify-content-start"
                                    )
                                )
                                .append(
                                  $("<p>")
                                    .addClass("card-text")
                                    .text(data[i].body)
                                )
                            )
                        );
                      } else {
                        $("#posts").append(
                          $("<div>")
                            .addClass("card")
                            .append(
                              $("<div>")
                                .addClass("card-body")
                                .append(
                                  $("<a>")
                                    .addClass("d-flex justify-content-start")
                                    .attr("id", "poster-data")
                                    .attr(
                                      "href",
                                      profile_base_path +
                                        data[i].poster_username
                                    )
                                    .append(
                                      $("<img>")
                                        .attr(
                                          "src",
                                          base_path + data[i].poster_avatar
                                        )
                                        .addClass("rounded-circle")
                                        .attr("height", "30")
                                        .attr("width", "30")
                                        .attr("loading", "lazy")
                                    )
                                    .append(
                                      $("<h5>")
                                        .text(data[i].poster_username)
                                        .attr("id", "post-name")
                                    )
                                    .attr("id", post_id)
                                )

                                .append(
                                  $("<h6>")
                                    .text(formattedDate)
                                    .addClass(
                                      "card-subtitle mb-2 text-muted d-flex justify-content-start"
                                    )
                                )
                                .append(
                                  $("<p>")
                                    .addClass("card-text")
                                    .text(data[i].body)
                                )
                            )
                            .append(
                              $("<img>")
                                .attr("src", base_path + img_src)
                                .addClass("card-img-top")
                            )
                        );
                      }
                      if (data[i].membership_locked) {
                        lockPath =
                          "http://localhost:9080/img/post-lock-icon.png";
                        $("#" + post_id).append(
                          $("<img>")
                            .attr("src", lockPath)
                            .attr("id", "locked-post")
                        );
                      }
                    }
                  }
                },
                error: function (jqXHR) {
                  console.log(jqXHR);
                  if (jqXHR.status === 403 || jqXHR.status === 401) {
                    alert("403 or 401");
                    window.location.replace("http://localhost:9080/login.html");
                  }
                },
              });
            } else {
              $.ajax({
                type: "POST",
                url: userPostsEndpoint,
                data: userPostsRequest,
                headers: {
                  Authorization: "Bearer " + $.cookie("access_token"),
                },
                success: function (data) {
                  console.log(data);
                  if (data == null) {
                    $("#posts").append(
                      $("<p>").text("У пользователя еще нет записей")
                    );
                  } else {
                    for (let i = 0; i < data.length; i++) {
                      let post_id = "post_" + i;
                      let img_src = data[i].image_ref;
                      posted_at = data[i].posted_at;
                      date = new Date(posted_at);
                      const formattedDate = date.toLocaleString("ru", {
                        day: "numeric",
                        month: "short",
                        year: "numeric",
                        hour: "numeric",
                        minute: "2-digit",
                      });
                      if (img_src === undefined) {
                        $("#posts").append(
                          $("<div>")
                            .addClass("card")
                            .append(
                              $("<div>")
                                .addClass("card-body")
                                .append(
                                  $("<a>")
                                    .addClass("d-flex justify-content-start")
                                    .attr("id", "poster-data")
                                    .attr(
                                      "href",
                                      profile_base_path +
                                        data[i].poster_username
                                    )
                                    .append(
                                      $("<img>")
                                        .attr(
                                          "src",
                                          base_path + data[i].poster_avatar
                                        )
                                        .addClass("rounded-circle")
                                        .attr("height", "30")
                                        .attr("width", "30")
                                        .attr("loading", "lazy")
                                    )
                                    .append(
                                      $("<h5>")
                                        .text(data[i].poster_username)
                                        .attr("id", "post-name")
                                    )
                                    .attr("id", post_id)
                                )
                                .append(
                                  $("<h6>")
                                    .text(formattedDate)
                                    .addClass(
                                      "card-subtitle mb-2 text-muted d-flex justify-content-start"
                                    )
                                )
                                .append(
                                  $("<p>")
                                    .addClass("card-text")
                                    .text(data[i].body)
                                )
                            )
                        );
                      } else {
                        $("#posts").append(
                          $("<div>")
                            .addClass("card")
                            .append(
                              $("<div>")
                                .addClass("card-body")
                                .append(
                                  $("<a>")
                                    .addClass("d-flex justify-content-start")
                                    .attr("id", "poster-data")
                                    .attr(
                                      "href",
                                      profile_base_path +
                                        data[i].poster_username
                                    )
                                    .append(
                                      $("<img>")
                                        .attr(
                                          "src",
                                          base_path + data[i].poster_avatar
                                        )
                                        .addClass("rounded-circle")
                                        .attr("height", "30")
                                        .attr("width", "30")
                                        .attr("loading", "lazy")
                                    )
                                    .append(
                                      $("<h5>")
                                        .text(data[i].poster_username)
                                        .attr("id", "post-name")
                                    )
                                    .attr("id", post_id)
                                )

                                .append(
                                  $("<h6>")
                                    .text(formattedDate)
                                    .addClass(
                                      "card-subtitle mb-2 text-muted d-flex justify-content-start"
                                    )
                                )
                                .append(
                                  $("<p>")
                                    .addClass("card-text")
                                    .text(data[i].body)
                                )
                            )
                            .append(
                              $("<img>")
                                .attr("src", base_path + img_src)
                                .addClass("card-img-top")
                            )
                        );
                      }
                      if (data[i].membership_locked) {
                        lockPath =
                          "http://localhost:9080/img/post-lock-icon.png";
                        $("#" + post_id).append(
                          $("<img>")
                            .attr("src", lockPath)
                            .attr("id", "locked-post")
                        );
                      }
                    }
                  }
                },
                error: function (jqXHR) {
                  console.log(jqXHR);
                  if (jqXHR.status === 403 || jqXHR.status === 401) {
                    alert("403 or 401");
                    window.location.replace("http://localhost:9080/login.html");
                  }
                },
              });
            }
          },
          function (value) {}
        );
      },
      error: function (jqXHR) {
        window.location.replace("http://localhost:9080/feed.html");
      },
    });
  },
  error: function (jqXHR) {},
});

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
      error: function() {
        window.location.replace("http://localhost:9080/feed.html");
      }
    });
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

async function getUserSubbedTiers(id) {
  return new Promise((resolve, reject) => {
    let userTiersPath = "http://localhost:8080/usertiers";
    $.ajax({
      type: "GET",
      url: userTiersPath,
      headers: {
        Authorization: "Bearer " + $.cookie("access_token"),
      },
      success: function (data) {
        resolve(data);
      },
      error: function () {
        resolve(false);
      },
    });
  }).catch();
}

function buildTierDivs(
  membership_data,
  userHasMembership,
  userTiers,
  loggedInUserID
) {
  $("#membership-info").append(
    $("<div>")
      .addClass("card-header")
      .attr("id", "membership-tiers")
      .text("Уровни платной подписки")
  );
  if (userHasMembership) {
    console.log(membership_data.tiers.length);
    for (let i = 0; i < membership_data.tiers.length; i++) {
      $("#membership-tiers").append(
        $("<li>")
          .addClass("list-group-item px-3 tier")
          .attr("id", "tier" + i)
          .append(
            $("<h5>", "")
              .attr("id", "tier-name")
              .text(membership_data.tiers[i].name)
          )
          .append(
            $("<h6>")
              .attr("id", "tier-name")
              .addClass("text-muted")
              .text(membership_data.tiers[i].price + "₽ в месяц")
          )
      );
      if (membership_data.tiers[i].image_ref !== undefined) {
        imgBasePath = "http://localhost:9080/img/";
        $("#tier" + i).append(
          $("<img>")
            .attr("id", "tier-img")
            .attr("src", imgBasePath + membership_data.tiers[i].image_ref)
        );
      }
      $("#tier" + i).append(
        $("<p>").addClass("card-text").text(membership_data.tiers[i].rewards)
      );
      if (membership_data.owner_id === loggedInUserID) {
        $("#tier" + i)
          .append($("<p>").attr("id", "subbed-label").text("Вы подписаны"))
          .append(
            $("<h6>")
              .attr("id", "tier-name")
              .addClass("text-muted")
              .text("Навсегда")
          );
        continue;
      }
      console.log(userTiers);
      if (userTiers !== undefined) {
        console.log("USER TIERS - ", userTiers.length);

        let userSubbed = false;
        let currentTier;
        for (let index = 0; index < userTiers.length; index++) {
          console.log("INDEX - ", index);
          console.log(
            userTiers[index].tier_id,
            membership_data.tiers[i].id,
            userTiers[index] == membership_data.tiers[i].id
          );
          if (userTiers[index].tier_id == membership_data.tiers[i].id) {
            userSubbed = true;
            currentTier = userTiers[index];
          }
        }
        if (userSubbed) {
          date = new Date(currentTier.member_until);
          const formattedDate = date.toLocaleString("ru", {
            day: "numeric",
            month: "short",
            year: "numeric",
          });

          $("#tier" + i)
            .append($("<p>").text("Вы подписаны").attr("id", "subbed-label"))
            .append(
              $("<h6>")
                .attr("id", "tier-name")
                .addClass("text-muted")
                .text("До " + formattedDate)
            );
        } else {
          $("#tier" + i).append(
            $("<a>")
            .addClass("btn btn-primary")
            .text("Подписаться")
            .attr("id","tier-btn"+i)
          );
          btn_id = "#tier-btn"+i
          tier_id = membership_data.tiers[i].id
          console.log(btn_id,tier_id)
          tierSubButtonListener(btn_id,tier_id)
        }
      } else {
        $("#tier" + i).append(
          $("<a>").addClass("btn btn-primary").text("Подписаться")
        );
        btn_id = "#tier-btn"+i
        tier_id = membership_data.tiers[i].id
        console.log(btn_id,tier_id)
        tierSubButtonListener(btn_id,tier_id)
      }
    }
  } else {
    $("#membership-tiers").append(
      $("<h6>")
        .attr("id", "tier-name")
        .addClass("text-muted")
        .text("У пользователя пока нет платных уровней подписки.")
    );
  }
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

async function getUserShop() {
  shopUserID = await getUserIDByUsername(username);
  loggedInUserID = await getLoggedInUserID();
  shopData = await getUserShopByOwnerID(shopUserID);
  buildShopDiv(shopData)
  console.log(shopData)
}

function buildShopDiv(shopData) {
  baseImgPath = "http://localhost:9080/img/"
  shopPath = "http://localhost:9080/user_shop.html?user="

  if (shopData !== null) {
    $("#shop-info")
    .append(
        $("<div>")
          .addClass("card-header")
          .text("Магазин")
          .attr("id","shop-showcase")
    )
    .append(
      $("<img>")
          .addClass("card-img-top")
          .attr("src", baseImgPath+shopData[0].image_ref)
    )
    .append(
      $("<h5>")
          .addClass("card-body")
          .attr("id","shop-showcase-name")
          .text(shopData[0].name)
    )
    .append(
      $("<h6>")
        .addClass("text-muted")
        .text(shopData[0].price+"₽")
    )
    .append(
      $("<a>").addClass("btn btn-primary").text("В магазин").attr("id","toShopBtn")
      .attr("href",shopPath+username)
    )
  } else {
    $("#shop-info")
    .append(
      $("<h5>")
        .text("Магазин")
        .attr("id","shop-showcase")
  )
    .append(
      $("<h6>")
      .addClass("text-muted")
      .text("У пользователя пока нет товаров.")
    )
  }


}

getUserShop();

async function getUserTiers() {
  loggedInUserID = await getLoggedInUserID();
  owner_id = await getUserIDByUsername(username);
  membership_id = await getMembershipIDByOwnerID(owner_id);
  if ($.cookie("access_token") !== undefined) {
    uuid = parseJwt($.cookie("access_token"));
    subbedTiers = await getUserSubbedTiers(uuid.sub);
    if (membership_id === false) {
      let data;
      buildTierDivs(data, false, subbedTiers, loggedInUserID);
    } else {
      membership_data = await getMembershipDataByID(membership_id);
      buildTierDivs(membership_data, true, subbedTiers, loggedInUserID);
    }
  } else {
    if (membership_id === false) {
      let data;
      buildTierDivs(data, false, false, loggedInUserID);
    } else {
      membership_data = await getMembershipDataByID(membership_id);
      buildTierDivs(membership_data, true, false, loggedInUserID);
    }
  }

}

getUserTiers();

function tierSubButtonListener(btnID, tierID) {

  subscribePath = "http://localhost:8080/tiersubscribe";
  console.log("listening to " + btnID)
  $(btnID).on("click", function () {

    userUUID = parseJwt($.cookie("access_token"));
    subTierJSON =
      '{"user_id":"' + userUUID.sub + '","tier_id":"' + tierID + '"}';
      console.log(subTierJSON)
    $.ajax({
      type: "POST",
      url: subscribePath,
      data: subTierJSON,
      headers: {
        Authorization: "Bearer " + $.cookie("access_token"),
      },
      success: function (data) {
        alert("ALL SUBS ARE FREE NOW FOR TESTING PURPOSES");
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

function subscribeButtonClickListener() {
  $("#sub_button").off();

  subscribePath = "http://localhost:8080/subscribe";
  unsubscribePath = "http://localhost:8080/unsubscribe";

  $("#sub_button").on("click", function () {
    unsubJSON = '{"id":"' + userPageID + '"}';
    if (subbed) {
      $.ajax({
        type: "POST",
        url: unsubscribePath,
        data: unsubJSON,
        headers: {
          Authorization: "Bearer " + $.cookie("access_token"),
        },
        success: function (data) {
          changeSubButton();
        },
        error: function (jqXHR) {
          if (jqXHR.status === 403 || jqXHR.status === 401) {
            alert("403 or 401");
            window.location.replace("http://localhost:9080/login.html");
          }
        },
      });
    } else {
      $.ajax({
        type: "POST",
        url: subscribePath,
        data: unsubJSON,
        headers: {
          Authorization: "Bearer " + $.cookie("access_token"),
        },
        success: function (data) {
          changeSubButton();
        },
        error: function (jqXHR) {
          if (jqXHR.status === 403 || jqXHR.status === 401) {
            alert("403 or 401");
            window.location.replace("http://localhost:9080/login.html");
          }
        },
      });
    }
  });
}
