const logos = [
  { src: "/icon/adidas.png", alt: "Adidas" },
  { src: "/icon/apple.png", alt: "Apple" },
  { src: "/icon/github.png", alt: "GitHub" },
  { src: "/icon/firefox.png", alt: "FireFox" },
  { src: "/icon/java-script.png", alt: "Java Script" },
  { src: "/icon/linkedin.png", alt: "Adidas" },
  { src: "/icon/nvidia.png", alt: "Apple" },
  { src: "/icon/tiktok.png", alt: "GitHub" },
  { src: "/icon/twitch.png", alt: "Fire Fox" },
];
export default function BrandIcons() {
  return (
    <div className="overflow-hidden w-full py-8 [mask-image:linear-gradient(to_right,transparent,black_10%,black_90%,transparent)]">
      <div className="flex w-max gap-16 items-center animate-marquee hover:[animation-play-state:paused]">
        {/* Original */}
        {logos.map((logo) => (
          <img
            key={logo.alt}
            src={logo.src}
            alt={logo.alt}
            className="h-24 grayscale opacity-50 hover:opacity-100 hover:grayscale-0 transition-all duration-300"
          />
        ))}

        {/* Clone — for seamless loop */}
        {logos.map((logo) => (
          <img
            key={`${logo.alt}-clone`}
            src={logo.src}
            alt={logo.alt}
            aria-hidden="true"
            className="h-24 grayscale opacity-50"
          />
        ))}
      </div>
    </div>
  );
}
