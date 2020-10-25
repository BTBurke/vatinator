module.exports = {
  future: {
    removeDeprecatedGapUtilities: true,
  },
  important: true,
  purge: ['./components/**/*.{js,ts,jsx,tsx}', './pages/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        'base': '#211F30',
        'secondary': '#32384D',
        'accent-1': '#E29930',
        'accent-2': '#217CA3',
      },
    },
    maxHeight: {
      '0': '0',
      '1/4': '25%',
      '1/4vh': '25vh',
      '1/3': '33%',
      '1/3vh': '33vh',
      '1/2': '50%',
      '2/3': '66%',
      '2/3vh': '66vh',
      '3/4': '75%',
      '3/4vh': '75vh',
      'full': '100%',
    },
  },
  variants: {},
  plugins: [
    require('tailwindcss-tables')({
      cellPadding: '.75rem',  // default: .75rem
      tableBorderColor: '#dee2e6',  // default: #dee2e6
      tableStripedBackgroundColor: 'rgba(0,0,0,.05)',  // default: rgba(0,0,0,.05)
      tableHoverBackgroundColor: 'rgba(0,0,0,.075)',  // default: rgba(0,0,0,.075)
      tableBodyBorder: true, // default: true. If set to false, borders for the table body will be removed. Only works for normal tables (i.e. does not apply to .table-bordered)
      verticalAlign: 'top', // default: 'top'
    }),
    require('@tailwindcss/custom-forms'),
  ],
}
