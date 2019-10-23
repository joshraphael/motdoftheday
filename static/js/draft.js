$(document).ready(function () {
    $(".history").hide()
    var init_id = $("#history").children(":selected").attr("id")
    $("#history-" + init_id).show()
    $("#history").change(function () {
        var id = $(this).children(":selected").attr("id")
        $(".history").hide();
        $("#history-" + id).show()

    });
    $("#edit-button").on("click", function () {
        var post_history_id = $("#history").children(":selected").attr("id");
        window.location.href = '/edit/' + post_history_id;
    })
})