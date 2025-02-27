<template>
  <router-view/>
</template>

<style>
:root {
  --primary-color: #3490dc;
  --secondary-color: #38c172;
  --dark-color: #2d3748;
  --light-color: #f8fafc;
  --gray-color: #e2e8f0;
  --danger-color: #e3342f;
  --success-color: #38c172;
  --border-radius: 8px;
  --box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  --transition: all 0.3s ease;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 1.6;
  color: var(--dark-color);
  background-color: #f0f4f8;
  margin: 0;
  padding: 0;
  min-height: 100vh;
}

#app {
  width: 100%;
  display: flex;
  justify-content: center;
  padding: 2rem 1rem;
}

@media (max-width: 768px) {
  #app {
    padding: 1rem 0.5rem;
  }
}

h1, h2, h3, h4, h5, h6 {
  margin-bottom: 1rem;
  font-weight: 600;
  line-height: 1.2;
}

a {
  color: var(--primary-color);
  text-decoration: none;
  transition: var(--transition);
}

a:hover {
  text-decoration: underline;
}

p {
  margin-bottom: 1rem;
}

/* For mobile */
@media (max-width: 640px) {
  :root {
    font-size: 14px;
  }
}
</style>