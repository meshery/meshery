import styled from "styled-components";

const NavigationWrap = styled.header`

  background-color:  "black";
  position: sticky;
  width: 100%;
  z-index: 9999;
  top: 0;
  transition: 0.8s cubic-bezier(0.2, 0.8, 0.2, 1);

  .toolbarDrawerClosed {
   

    background-color: #fff;
  }
  .toolbarDrawerOpen {

    background-color: #000;
  }

  .appBarOnDrawerOpen {
  
    background-color: #000;
  }

  .nav-container {
    display: flex;
    @media (min-width:912px) and (max-width:992px){
      max-width:850px;
    }
  }
  .meshery-cta {
    position: relative;
    display: flex;
    align-items: center;
    right: -4.5rem;
    top: 1.5rem;
    max-height: 50px;

    a {
      display: flex;
    }

    @media screen and (max-width: 1550px) {
      position: relative;
      right: -4rem;
      top: 1.5rem;
    }
    @media screen and (max-width: 1275px) {
      right: -1.95rem;
    }
    @media screen and (max-width: 992px) {
      right: -1rem;
    }
    @media screen and (max-width: 912px) {
      a {
        display: none;
      }
      .dark-theme-toggle {
        margin-left: 0.5em;
        margin-bottom: 2.5em;
      }
 
    }
  }
  .navbar-wrap {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    position: relative;
  }
  .mobile-menu-icon {
    display: none;
    color: ${props => props.theme.menuColor};
    font-size: 24px;
  }
  .dropdown {
    position: absolute;
    background: ${props => props.theme.body};
    opacity: 0;
    border: 1px solid ${props => props.theme.DarkTheme ? "rgb(20, 20, 20)" : "#f5f5f5"};    
    border-radius: 50px;
    top: 100%;
    left: 0;
    width: 120%;
    display: block;
    box-shadow: 0px 5px 10px 1px rgba(0, 179, 159, 0.50);
    animation: bobbleout ease .18s forwards;
    pointer-events: none;
  }
  .wrap {
    display: block;
    position: absolute;
    top: 100%;
    left: 0;
    width: 100%;
    height: 350px;
    margin-top: -2rem; 
    visibility: visible;
  }
  .mobile-dropdown {
    display: none;
  }
  .title {
    color: black;
    font-size: 20px;
    font-weight: 600;
  }
  .arrow-icon {
    background: #00B39F;
    border-radius: 25%;
  }
  .nav-grid {
    display: grid;
    grid-template-columns: 35% 65%;
    .hr {
      background: ${props => props.theme.DarkTheme ? "rgb(20, 20, 20)" : "rgb(250, 250, 250)"};
      display: flex;
      flex-direction: column;
      position: relative;
      padding: 3em;
      border-radius: 50px 0 0 50px;
      .sub-item {
        padding: 0;
      }
      .section {
        color: ${props => props.theme.text};
        font-size: 20px;
        font-weight: 600;
        margin-left:0px;
        padding-left:0px;
      }
      li:nth-last-child(2) {
        margin-bottom: auto;
      }
      .action-items {
        margin-top: 10px;
        position: relative;
        .action-link {
          padding: 2px 0;
          display: block;
          .readmore-btn {
            color: ${props => props.theme.text};
          }
          &:before {
            content: none;
          }
          &:hover{
            .readmore-btn {
                color: ${props => props.theme.menuHoverColor};
                svg{
                    margin-left: 3px;
                      transform: scale(1.2);
                }
            }
          }
        }
      }
    }
    .nav-display {
      border-left: 2px solid ${props => props.theme.DarkTheme ? "rgb(60, 60, 60)" : "#f1f1f1"}; 
      padding-top: 1em;
      border-radius: 0 50px 50px 0;
      display: grid;
      grid-template-columns: 50% 50%;
      @media screen and (max-width: 992px) and (min-width: 912px) {
         padding: .5em;
       }
    }
  }
  @keyframes nav-bar-anim {
    0% {
      opacity: 0;
      transform: translate3d(0px, 0px, 0px) scale(.5, .5);
    }

    100% {
      transform: translate3d(0px, 0px, 200px) scale(1, 1);
      opacity: 1;
    }
  }
  @keyframes bobble {
    0% {
        transform: translate3d(0px, -15px, 0px) scale(.75, .75); opacity: 0;
    }
    50% {
        transform: translate3d(0px, 0px, 4px) perspective(800px) ;
        opacity: 70%
    }
    100% {
        transform: translate3d(0px, 0px, 8px) scale(1, 1);
        opacity: 100%
    }
  }
  @keyframes bobbleout {
    0% {
        transform: translate3d(0px, 0px, 8px) scale(.9, .9);
        opacity: 40%;
    }
    50% {
        transform: translate3d(0px, 0px, 4px) scale(.75, .75);
        perspective(800px);
        opacity: 20%;
    }
    100% {
        transform: translate3d(0px, -15px, 0px) scale(.5, .5);
        opacity: 0;
    }
  }
  @keyframes flowin {
    0% {
      transform: translate3d(-10px, 0px, 0px);
      opacity: 0;
    }
    50% {
      transform: translate3d(10px, 0px, 0px);
      opacity: 0.5;
    }
    100% {
      transform: translate3d(0px, 0px, 0px);
      opacity: 1;
    }
  }
  ul:hover > ul {
    animation: bobble ease .3s forwards;
    pointer-events: auto;
    visibility: visible;
  }
  nav {
    ul {
      margin: 0;
      padding: 0;
      list-style: none;
      display: flex;
    }
    ul.collapsed {
      margin: 2rem 0 2rem 0;
    }
    .nav-item {
      position: relative;
      overflow: hidden;
      &:hover {
        .menu-link {
          span {
            color: ${props => props.theme.menuHoverColor};
            &:before {
              opacity: 1;
              animation: flowin .5s ease-in;
            }
          }
        }
      }
    }
    .dropdown {
      li {
        padding-top: 5px;
        a {
          display: block;
          &:hover {
            color: ${props => props.theme.menuHoverColor}; !important;
          }
          &:before {
            content: none;
          }
        }
      }
    }
    a,
    .nav-active {
      position: relative;
      color: ${props => props.theme.menuColor};
      display: flex;
      line-height: 1.5rem;
      font-size: 15px;
      transition: 450ms all;
      padding: 0px 20px 5px 20px;
      cursor: pointer;
      &:before {
        content: "";
        position: absolute;
        left: 20;
        bottom: 0;
        width: 20px;
        height: 1px;
        opacity: 0;
        background: ${props => props.theme.menuHoverColor};
        transition: 450ms all;
      }
      &:hover {
        color: ${props => props.theme.menuHoverColor};
        &:before {
          opacity: 1;
        }
      }
    }
    .menu-link {
      margin: auto 1.25rem;
      padding: 0.25rem 0 0.25rem;
      &:before {
        content: none;
      }
      span {
        &:before {
          content: "";
          position: absolute;
          bottom: 0;
          width: 100%;
          height: 1px;
          opacity: 0;
          background: ${props => props.theme.menuHoverColor};
        }
      }
    }
  }
  .logo {
    margin-top: 8px;
    transition: all 0.8s cubic-bezier(0.3, 0.8, 0.2, 1) 0s;
    img {
      width: 155px;
    }
    @media screen and (max-width: 992px) and (min-width: 912px) {
      width: 125px;
    }
  }
  &.scrolled {    
    box-shadow: rgba(0, 179, 159, 0.2) 0px 10px 25px;
    .nav {
      .nav-item{
        a {
          color:${props => props.theme.menuColor};
          &:hover {
            color: ${props => props.theme.menuHoverColor};
          }
        }
      }
    }
    .wrap {
      height: 300px;
      margin-top: -34px;
    }
  }

  @media only screen and (max-width: 912px) {
    height: auto;
    min-height: 50px;
    padding: 15px 20px;
    .nav {
      margin-top: 7px;
    }
    &.scrolled {
      padding: 15px 20px;
    }
    ul.collapsed {
      display: none;
    }
    .navbar-wrap {
      position: static;
      display: flex;
      align-items: center;
    }
    .logo {
      width: 100%;
      display: block;
      margin-top: 7px;
      img {
        width: 130px;
      }
    }
    .mobile-dropdown-container {
      position: absolute;
      top: auto;
      left: 0;
      height: auto;
      width: 100%;
      padding: 0 10px;
      visibility: hidden;
      opacity: 0;
      transition: .3s;
    }
    .mobile-dropdown {
      position: relative;
      padding: 10px 0 10px 15px;
      display: block;
      background: ${props => props.theme.DarkTheme ? "rgb(20, 20, 20)" : "rgb(250, 250, 250)"};
      border-radius: 10px;
      box-shadow: 0px 10px 10px 10px rgba(0, 211, 169, 0.10);
      max-height: 400px;
      overflow-y: scroll;
    }
    .expand {
      visibility: visible;
      opacity: 1;
      .mobile-dropdown {
        visibility: visible;
        opacity: 1;
      }
    }

    .mobile-collapsed {
      display: flex;
      flex-direction: column;
    }
    .mobile-menu-icon {
      display: inline-block;
      position: relative;
      cursor: pointer;
      animation: open-icon 0.3s ease-in;
    }
    .mobile-menu-icon.open {
      animation: close-icon 0.3s ease-in;
    }
    @keyframes open-icon {
      from { opacity: 0; }
      to { opacity: 1; }
      0% { transform: scale(0); }
      100% { transform: scale(1); }
    }
    @keyframes close-icon {
        from { opacity: 0; }
        to { opacity: 1; }
        0% { transform: rotate(0deg); }
        100% { transform: rotate(180deg); }
    }
    .mobile-nav-item {
      padding: 1px;
      .menu-item {
        font-size: 16px;
        font-weight: 600;
      }
      a:before {
        content: none;
      }
      ul {
        display: list-item;
      }
    }
    .mobile-nav-subitem {
      padding-left: 10px;
    }
  }
  
  .nav-link-active {
    color: ${(props) => props.theme.menuHoverColor};
  }
  
  .anchor:before {
    content: none;
  }
  .post-block{
    display: flex;
    flex-direction: column;
    align-items: stretch;
    min-height: 300px;
    overflow: hidden;
    margin: 0 auto;
    padding: 1rem 1rem 0rem 1rem;
    &:hover{
      box-shadow: 0px 0px 8px -2px rgba(0, 0, 0, 0.25);
      border-radius: 1rem;
      .readmore-btn{
          color: ${props => props.theme.menuHoverColor};
          svg{
              margin-left: 3px;
                transform: scale(1.2);
          }
      }
      .post-title{
          color: ${props => props.theme.text};
      }
    }
    a:before {
      content: none;
    }
  }

  .readmore-btn-wrapper{
    display: flex;
    justify-content: space-between;
  }

  .post-thumb-block{
    overflow: hidden;
    border-radius: 10px;  
    height: 160px;
    max-width: 100%;
    margin: 0 auto;
    .gatsby-image-wrapper, .old-gatsby-image-wrapper{
        height:100%;
        transition: all 0.3s ease-in;
    }
    img{
        max-width: 240px;
        max-height: 160px;
        width: auto;
        height: inherit;
        display: block;
        text-align: center;
        margin-left: auto;
    }
  }
  .post-content-block{
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    justify-content: space-between;
    border-radius: 0 0 3px 3px;
    position: relative;
    padding: 20px 0px;
  }
  .post-title{
    font-size: 16px;
    font-weight: 600;
    line-height: 28px;
    color: ${props => props.theme.lightTextColor};
    margin-bottom: 20px;
    display: block;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    -webkit-transition: 450ms all;
    transition: 450ms all;
    &:hover{
      color: ${props => props.theme.highlightLightColor};


    }
    
    @supports (-webkit-line-clamp: 2) {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: initial;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }
  .readmore-btn {
    color: ${props => props.theme.lightTextColor};
    display: flex;
    align-items: center;
    transition: all 0.3s linear;
    padding-left: 0px;
    svg{
      margin-left: 0px;
      font-size: 27px;
      transition: all 0.3s linear;
    }
  }
  .dark-theme-toggle {

    margin-left: 2em;


  }

  .toggle {
    --size: 1.5rem;
    appearance: none;
    outline: none;
    cursor: pointer;
    width: var(--size);
    height: var(--size);
    box-shadow: inset calc(var(--size) * 0.33) calc(var(--size) * -0.25) 0;
    border-radius: 999px;
    color: #00B39F;
    transition: all 500ms;
  }

  .toggle:checked {
    --ray-size: calc(var(--size) * -0.4);
    --offset-orthogonal: calc(var(--size) * 0.65);
    --offset-diagonal: calc(var(--size) * 0.45);
    transform: scale(0.75);
    color: #3c494f;
    box-shadow: inset 0 0 0 var(--size), calc(var(--offset-orthogonal) * -1) 0 0 var(--ray-size), var(--offset-orthogonal) 0 0 var(--ray-size), 0 calc(var(--offset-orthogonal) * -1) 0 var(--ray-size), 0 var(--offset-orthogonal) 0 var(--ray-size), calc(var(--offset-diagonal) * -1) calc(var(--offset-diagonal) * -1) 0 var(--ray-size), var(--offset-diagonal) var(--offset-diagonal) 0 var(--ray-size), calc(var(--offset-diagonal) * -1) var(--offset-diagonal) 0 var(--ray-size), var(--offset-diagonal) calc(var(--offset-diagonal) * -1) 0 var(--ray-size);
  }

  .toggle {
    z-index: 1;
  }

  .toggle:checked~.background {
    --bg: white;
  }

  .toggle:checked~.title {
    --color: #fa0;
  }
  
`;

export default NavigationWrap;
