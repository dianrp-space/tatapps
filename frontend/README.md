# TatApps Frontend

Frontend aplikasi TatApps menggunakan Vue 3, Vite, dan TailwindCSS.

## Quick Start

```bash
# Install dependencies
npm install
# atau
yarn install

# Setup environment
cp .env.example .env

# Run development server
npm run dev
# atau
yarn dev
```

Frontend akan berjalan di `http://localhost:5173`

## Build for Production

```bash
npm run build
# atau
yarn build

# Preview production build
npm run preview
# atau
yarn preview
```

## Project Structure

```
src/
├── api/           # API client & axios configuration
├── assets/        # Static assets (images, fonts, etc)
├── components/    # Reusable Vue components
├── layouts/       # Layout components
├── router/        # Vue Router configuration
├── stores/        # Pinia stores (state management)
├── views/         # Page components
├── App.vue        # Root component
└── main.js        # Application entry point
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build

## Technologies

- Vue 3 (Composition API)
- Vite
- Vue Router
- Pinia (State Management)
- Axios (HTTP Client)
- TailwindCSS (Styling)
- PrimeVue (UI Components)
