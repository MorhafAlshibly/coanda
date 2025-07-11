// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const { themes } = require("prism-react-renderer");
const lightCodeTheme = themes.github,
	darkCodeTheme = themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
	title: "coanda",
	tagline: "A backend game management system",
	url: "https://morhafalshibly.github.io",
	baseUrl: "/coanda",
	onBrokenLinks: "throw",
	onBrokenMarkdownLinks: "warn",
	favicon: "img/favicon.ico",
	organizationName: "MorhafAlshibly", // Usually your GitHub org/user name.
	projectName: "coanda", // Usually your repo name.
	plugins: ["@graphql-markdown/docusaurus"], // See .graphqlrc for configuration
	presets: [
		[
			"classic",
			/** @type {import('@docusaurus/preset-classic').Options} */
			({
				blog: false,
				docs: {
					routeBasePath: "/",
				},
				theme: {
					customCss: "./src/css/custom.css",
				},
			}),
		],
	],

	themeConfig:
		/** @type {import('@docusaurus/preset-classic').ThemeConfig} */
		({
			navbar: {
				title: "coanda",
				logo: {
					alt: "wind-solid",
					src: "img/wind-solid.svg",
				},
				items: [
					{
						href: "https://github.com/MorhafAlshibly/coanda",
						label: "GitHub",
						position: "right",
					},
				],
			},
			footer: {
				style: "light",
				links: [],
				copyright: `Copyright © ${new Date().getFullYear()} coanda. Built with GraphQL-Markdown & Docusaurus.`,
			},
			prism: {
				theme: lightCodeTheme,
				darkTheme: darkCodeTheme,
			},
		}),
};

module.exports = config;
