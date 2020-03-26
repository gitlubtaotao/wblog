$(document).ready(function () {
    $(".read_comment").on("click", function (e) {
        $.post($(e.target).data("href"), {}, function (result) {
            window.location.href = $(e.target).data("redirect");
        }, 'json');
    });

    $(".read_all").on("click", function (e) {
        $.post("/admin/read_all", {}, function (result) {
            window.location.reload();
        }, "json");
    });
});