// Example starter JavaScript for disabling form submissions if there are invalid fields
    // source: https://getbootstrap.com/docs/4.1/examples/checkout/?
    (function() {
    'use strict';

    window.addEventListener('load', function() {
        // Fetch all the forms we want to apply custom Bootstrap validation styles to
        var forms = document.getElementsByClassName('needs-validation');

        // Loop over them and prevent submission
        var validation = Array.prototype.filter.call(forms, function(form) {
        form.addEventListener('submit', function(event) {
            if (form.checkValidity() === false) {
            event.preventDefault();
            event.stopPropagation();
            }
            form.classList.add('was-validated');
        }, false);
        });
    }, false);
    })();









$('#loadTestSuccess').hide();
$('#loadTestError').hide();

$('#istioConfigSuccess').hide();
$('#istioConfigError').hide();

$("#loadTestForm").on( "submit", function( event ) {
    event.preventDefault();
    $('#loadRun').addClass("disabled");
    $theForm = $(this)[0];
    if ((typeof($theForm.checkValidity) == "function") && !$theForm.checkValidity()) {
        return;
    }
    startCountDownCounter($theForm.t.value);

    data = $( this ).serialize();
    console.log("serialized data: "+data);
    jQuery.ajax({
        url: "/play/load-test",
        type: "POST", 
        data: data, 
        error: function(){
        $('#loadRun').removeClass('disabled');
        $("#loadTestError").fadeTo(2000, 500); 
        },
        success: function(result){
        $('#loadRun').removeClass('disabled');
        $("#loadTestSuccess").fadeTo(2000, 500).slideUp(500, function(){
            $("#loadTestSuccess").slideUp(500);
            $('#ldClose').click();
        }); 
        var res = JSON.parse(result);
        $('#loadTestRunTime').text(new Date(res.StartTime));
        var data = fortioResultToJsChartData(res);
        showChart(data);
        $(".loadTestResults").show();
        }
    });
    });

    $("#istioForm").on('submit', function( event ) {
    event.preventDefault();
    $theForm = $(this)[0];
    if ((typeof($theForm.checkValidity) == "function") && !$theForm.checkValidity()) {
        return;
    }
    data = $( this ).serialize();
    // console.log("serialized data: "+data);

    jQuery.ajax({
        url: "/play/mesh",
        type: "POST", 
        data: data, 
        error: function(){
        $("#istioConfigError").fadeTo(2000, 500); 
        },
        success: function(result){
        $("#istioConfigSuccess").fadeTo(2000, 500).slideUp(500, function(){
            $("#istioConfigSuccess").slideUp(500);
            $('#isClose').click();
        }); 
        }
    });
    });

    function resetCounter(){
        $('#countdown #hour').text('00');
        $('#countdown #mins').text('00');
        $('#countdown #secs').text('00');
    }

    function startCountDownCounter(startMin){
    resetCounter();
    if(startMin.length == 1){
        startMin = '0' + startMin;
    }
    $('#ltCountDown #mins').html(startMin);
    var startTime = new Date();
    
    var intv = setInterval(function time(){
        var d = startMin * 60 - Math.abs(new Date().getTime() - startTime.getTime()) / 1000;

        var days = Math.floor(d / 86400);
        d -= days * 86400;

        var hrs = Math.floor(d / 3600) % 24;
        d -= hrs * 3600;

        var mins = Math.floor(d / 60) % 60;
        d -= mins * 60;

        var secs = Math.floor(d % 60);

        var hr = hrs + '';
        if(hr.length == 1){
            hr = '0' + hrs;
        }

        var min = mins + '';
        if(min.length == 1){
            min = '0' + mins;
        }
        var sec = secs + '';
        if(sec.length == 1){
            sec = '0' + secs;
        }
        $('#ltCountDown #hour').html(hr);
        $('#ltCountDown #mins').html(min);
        $('#ltCountDown #secs').html(sec);

        if ($('#ltCountDown #mins').html() == "00" && $('#ltCountDown #secs').html() == "00") {
        clearInterval(intv);
        }

    }, 1000);
    }