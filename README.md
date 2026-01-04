# bhelper

Utilities web application.

## Directory Structure

```
bhelper/
├── index.html
├── src/
│   └── components/
│       └── {componentName}/
│           ├── {componentName}.html
│           └── {componentName}.js
└── README.md
```

## Component Structure

Each component is a self-contained feature with:

- `{componentName}.html` - Template markup with Alpine.js directives
- `{componentName}.js` - Logic exported as a named function following `{componentName}Logic` pattern

## Adding New Components

1. Create directory: `src/components/{componentName}/`
2. Add template file: `{componentName}.html`
3. Add logic file: `{componentName}.js`
4. Update `index.html` to import and load the component

## Technologies

- Alpine.js for reactive components
- Tailwind CSS for styling
