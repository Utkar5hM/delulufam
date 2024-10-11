/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "./views/*.templ",
    'node_modules/preline/dist/*.js',
  ],
  theme: {
    extend: {},
  },
	plugins: [require("@tailwindcss/typography"), 
    require('preline/plugin'),],
}

