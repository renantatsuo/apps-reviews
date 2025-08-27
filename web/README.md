# Apps Reviews Web Client

Built on top of Vite powered with typescript and react-query.

## Architecture

The application follows a clean, component-based architecture with clear separation of concerns:

**Application (`src/App.tsx`)**

**Components (`src/components/`)**

**Library code (`src/lib/`)**

## Features

- **App Management**: Browse and select from monitored apps with thumbnails and names
- **Add New Apps**: Add apps to monitor by Apple App Store ID through the web interface
- **Real-time Reviews**: Display recent reviews from the last 48 hours for selected apps
- **Rating Display**: Visual star ratings with numeric scores
- **Responsive Design**: Works seamlessly on desktop and mobile devices
- **Error Handling**: Graceful error states with retry functionality
- **Loading States**: Smooth loading indicators during data fetching
- **Data Caching**: Client-side caching with React Query for better performance

## Running the Application

### Prerequisites

- Node.js (v16 or higher)
- npm or yarn package manager
- Backend server running on port 8080 (see `../server/README.md`)

### Development Server

Install dependencies and start the development server:

```bash
npm install
npm run dev
```

The application will start at `http://localhost:5173` and automatically proxy API requests to the backend server at `http://localhost:8080`.

## API Integration

The web client connects to the backend server through a proxy configuration:

- **Development**: Vite dev server proxies `/api/*` requests to `http://localhost:8080`
- **Production**: Requires proper reverse proxy configuration (nginx, Apache, etc.)

### API Endpoints Used

- `GET /api/apps` - Fetch all monitored apps
- `GET /api/apps/{appID}` - Fetch details for a specific app
- `POST /api/apps/{appID}` - Add a new app to monitor
- `GET /api/reviews/{appID}` - Fetch reviews for a specific app

## Technology Stack

| Technology      | Version | Purpose                              |
| --------------- | ------- | ------------------------------------ |
| **React**       | ^19.1.1 | UI framework                         |
| **TypeScript**  | ~5.8.3  | Type safety and developer experience |
| **Vite**        | ^7.1.2  | Build tool and development server    |
| **React Query** | ^5.85.5 | Data state management and caching    |
| **ESLint**      | ^9.33.0 | Code linting and formatting          |

## Usage

1. **Start the Backend**: Ensure the Go server is running (see `../server/README.md`)
2. **Start the Frontend**: Run `npm run dev` to start the web application
3. **Browse Apps**: View all currently monitored apps with their thumbnails and names
4. **Add New Apps**: Use the "Add New App" form to add apps by their Apple App Store ID (e.g., `1458862350`)
5. **Select Apps**: Click on any app from the list to view its recent reviews
6. **View Reviews**: Browse through recent reviews with ratings and content

## Configuration

The application can be configured through environment variables or by modifying the Vite configuration:

### Development Proxy

The development server automatically proxies API requests to the backend. To change the backend URL, modify `vite.config.ts`:

```typescript
server: {
  proxy: {
    "/api": {
      target: "http://your-backend-url:8080",
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ""),
    },
  },
},
```
