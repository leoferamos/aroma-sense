import React from 'react';


const wordLines = [
  ['AROMA', 'SENSE'],
  ['FEMALE'],
  ['MALE', 'SCENT'],
  ['NOTES'],
  ['FLORAL', 'DESIGN'],
  ['LOVE'],
  ['SEXY', 'AROMA'],
  ['CITRUS'],
  ['MUSKY', 'PURE'],
  ['WOODY'],
  ['VIBE'],
];

const WordGrid: React.FC = () => {
  return (
    <div
      className="relative h-full w-full flex flex-col pt-8 gap-y-4 overflow-hidden pl-0"
      style={{ fontFamily: 'Playfair Display, serif', fontWeight: 600 }}
    >
  {/* Overlay gradient with blur above the words */}
      <div
        className="absolute left-0 right-0 bottom-0 top-24 pointer-events-none z-20"
        style={{
          background: 'linear-gradient(to bottom, rgba(255,255,255,0) 0%, rgba(200,200,200,0.2) 40%, rgba(200,200,200,0.5) 80%, rgba(200,200,200,0.9) 100%)',
          filter: 'blur(32px)',
        }}
      ></div>
  {/* Words */}
      <div className="relative z-10">
        {wordLines.map((line, lineIdx) => (
          <div key={lineIdx} className={`flex ${line.length === 2 ? 'flex-row gap-x-4' : 'justify-start'}`}>
            {line.map((word, idx) => (
              <span
                key={idx}
                className={`text-black font-bold tracking-tight select-none text-left ${line.length === 2 ? 'text-4xl sm:text-5xl md:text-6xl lg:text-7xl xl:text-8xl' : 'text-5xl sm:text-6xl md:text-7xl lg:text-8xl xl:text-9xl'}`}
                style={{ wordBreak: 'break-word', overflow: 'hidden' }}
              >
                {word}
              </span>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

export default WordGrid;
