# SummVideo Frontend

Vue.js frontend for the SummVideo application.

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Lints and fixes files
```
npm run lint
```

## Integration with Backend

The built app will be placed in the `/static/vue/` directory of the backend, and the backend will serve it at the root route. Make sure to run the build command before deploying the backend application.

## Configuration

The Vue.js application is configured to proxy API requests to the backend during development. See `vue.config.js` for more details.