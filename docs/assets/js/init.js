 (function($){
   $(function(){

     $('.button-collapse').sideNav();
     $('.parallax').parallax();
     $('.collapsible').collapsible();
     $('.carousel.carousel-slider').carousel({fullWidth: true});
     $('.materialboxed').materialbox();
     $('.scrollspy').scrollSpy();
     $('.tap-target').tapTarget('open');

     if (localStorage.getItem('cookieconsent') === 'true') {
       $('#cookies').hide()
     }

     jQuery('#cookies').on('click', function(event) {
            localStorage.setItem('cookieconsent', 'true')
            jQuery('#cookies').toggle('hide');
       });

   }); // end of document ready
 })(jQuery);
