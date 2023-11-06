<template>
  <div class="about">
    <h1>{{ msg }}</h1>
  </div>
  <Footer />
</template>

<script>
import Footer from '@/components/Footer.vue'
import axios from 'axios'

export default {
  name: 'HomeView',
  components: {
    Footer,
  },
  data() {
    return {
      msg: 'Welcome to VMQ-GO!',
    };
  }
};
</script>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: calc(100vh - 90px);
    display: flex;
    align-items: center;
  }

  .about h1 {
    font-size: 3rem;  /* 48px */
    color: #fff; /* #fff */
    text-align: center; /* center */
    width: 100%; /* 100% */
    text-shadow: 0 0 10px #000;
  }

  .footer .t-link {
    color: #fff;
    text-shadow: 0 0 10px #000;
  }

  body {
    background-image: url('https://cdn.astwy.com/img/2023/06/8f7c3ffe634273147114977b1c0239fd.webp');
    background-size: cover;
    background-repeat: no-repeat;
    background-position: center center;
  }
}
</style>
