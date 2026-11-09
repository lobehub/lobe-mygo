import { mygoFrameSrc } from '../components/mygoFrame';

export function ShowcaseFrame() {
  return (
    <div className="mygo-showcase">
      <img alt="" aria-hidden className="mygo-showcase-poster mygo-shot-light" src="/home/showcase-light.png" />
      <img alt="" aria-hidden className="mygo-showcase-poster mygo-shot-dark" src="/home/showcase-dark.png" />
      <iframe
        className="mygo-showcase-live"
        loading="lazy"
        src={mygoFrameSrc('home/showcase')}
        title="A chat app built with Lobe MyGo, running live through WebAssembly"
      />
    </div>
  );
}
