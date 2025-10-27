import React, { useEffect, useRef, useState } from 'react';

const WORDS = [
  'AROMA',
  'SENSE',
  'FEMALE',
  'MALE',
  'SCENT',
  'NOTES',
  'FLORAL',
  'DESIGN',
  'LOVE',
  'SEXY',
  'AROMA',
  'CITRUS',
  'MUSKY',
  'PURE',
  'WOODY',
  'VIBE',
  'FRAGANCE',
];

function measureTextWidth(text: string, fontSize: number, fontFamily = "Playfair Display", fontWeight = 600) {
  const canvas = measureTextWidth.canvas || (measureTextWidth.canvas = document.createElement('canvas'));
  const ctx = canvas.getContext('2d')!;
  ctx.font = `${fontWeight} ${fontSize}px ${fontFamily}`;
  return ctx.measureText(text).width;
}
measureTextWidth.canvas = undefined as unknown as HTMLCanvasElement | undefined;

const WordGrid: React.FC = () => {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const firstLineRef = useRef<HTMLDivElement | null>(null);
  const [lines, setLines] = useState<string[][]>([]);
  const [fontSizeState, setFontSizeState] = useState<number>(32);
  const [lineHeightState, setLineHeightState] = useState<number>(36);
  const [overlayWidth, setOverlayWidth] = useState<number | null>(null);
  const [overlayParentTop, setOverlayParentTop] = useState<number | null>(null);
  const [overlayParentHeight, setOverlayParentHeight] = useState<number | null>(null);
  const [startBlurAt, setStartBlurAt] = useState<number | null>(null);
  const [overlayAbsLeft, setOverlayAbsLeft] = useState<number | null>(null);
  const [overlayAbsTop, setOverlayAbsTop] = useState<number | null>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    let mounted = true;

    const compute = () => {
      const el = containerRef.current as unknown as Element | null;
      if (!el || !(el instanceof Element)) return;
      const style = getComputedStyle(el);
      const paddingLeft = parseFloat(style.paddingLeft || '0');
      const paddingRight = parseFloat(style.paddingRight || '0');
      const availableWidth = Math.max(0, el.clientWidth - paddingLeft - paddingRight);
      const availableHeight = Math.max(0, el.clientHeight);
      const packAt = (fontSize: number) => {
        const gap = Math.max(4, Math.round(fontSize * 0.12));
        const lineHeight = Math.round(fontSize * 1.08);
        const resultLines: string[][] = [];
        let currentLine: string[] = [];
        let currentWidth = 0;
        let totalHeight = 0;

        for (let i = 0; i < WORDS.length; i++) {
          const word = WORDS[i];
          const w = Math.ceil(measureTextWidth(word, fontSize));
          if (w > availableWidth) continue; // skip if word can't fit at this fontSize

          const addedWidth = currentLine.length === 0 ? w : currentWidth + gap + w;
          if (addedWidth <= availableWidth) {
            if (currentLine.length === 0) {
              currentWidth = w;
              currentLine.push(word);
            } else {
              currentWidth = currentWidth + gap + w;
              currentLine.push(word);
            }
          } else {
            if (currentLine.length > 0) {
              totalHeight += lineHeight;
              if (totalHeight > availableHeight) return { lines: resultLines, totalHeight, lineHeight };
              resultLines.push(currentLine);
            }
            currentLine = [word];
            currentWidth = w;
          }
        }
        if (currentLine.length > 0) {
          totalHeight += lineHeight;
          if (totalHeight <= availableHeight) resultLines.push(currentLine);
        }
        return { lines: resultLines, totalHeight, lineHeight };
      };

      const minFont = 28;
      const maxFont = Math.max(120, Math.round(availableWidth * 0.16));
      let low = minFont;
      let high = maxFont;
      let best = minFont;
      let bestPack = packAt(best);

      while (low <= high) {
        const mid = Math.floor((low + high) / 2);
        const pack = packAt(mid);
        if (pack.totalHeight > 0 && pack.totalHeight <= availableHeight && pack.lines.length > 0) {
          best = mid;
          bestPack = pack;
          low = mid + 1;
        } else {
          high = mid - 1;
        }
      }

      // apply best found sizes
      setFontSizeState(best);
      setLineHeightState(bestPack.lineHeight || Math.round(best * 1.08));
      if (mounted) setLines(bestPack.lines || []);
    };

    compute();

    const handleResize = () => {
      compute();
    };

  const ro = new ResizeObserver(handleResize);
  if (containerRef.current) ro.observe(containerRef.current);
    window.addEventListener('orientationchange', handleResize);

    return () => {
      mounted = false;
      ro.disconnect();
      window.removeEventListener('orientationchange', handleResize);
    };
  }, []);
  useEffect(() => {
    const measure = () => {
      if (!containerRef.current) return;
      const container = containerRef.current;
      const parent = container.parentElement;
      const panel = parent && parent.parentElement ? parent.parentElement : parent;
          if (panel) {
            const panelRect = panel.getBoundingClientRect();
            const containerRect = container.getBoundingClientRect();
            const buffer = 6;
            const width = Math.max(0, Math.round(panelRect.width) + buffer * 2);
            const topRel = Math.max(0, Math.round(panelRect.top - containerRect.top - buffer));
            const height = Math.max(0, Math.round(panelRect.height) + buffer * 2);
            setOverlayWidth(width);
            setOverlayParentTop(topRel);
            setOverlayParentHeight(height);
            // absolute viewport coords
            setOverlayAbsLeft(Math.round(panelRect.left) - buffer);
            setOverlayAbsTop(Math.round(panelRect.top) - buffer);
          } else {
            setOverlayWidth(null);
            setOverlayParentTop(null);
            setOverlayParentHeight(null);
            setOverlayAbsLeft(null);
            setOverlayAbsTop(null);
          }

      // compute where blur should start relative to parent panel
      if (!firstLineRef.current) {
        const fallback = Math.round(lineHeightState + 1);
        setStartBlurAt(fallback);
      } else {
        const firstRect = firstLineRef.current.getBoundingClientRect();
        const panelRect = panel ? panel.getBoundingClientRect() : container.getBoundingClientRect();
        const start = Math.max(0, firstRect.bottom - panelRect.top);
        setStartBlurAt(Math.round(start));
      }
    };

  
    requestAnimationFrame(measure);

    const ro = new ResizeObserver(() => requestAnimationFrame(measure));
    if (containerRef.current) ro.observe(containerRef.current);
    window.addEventListener('orientationchange', measure);

    return () => {
      ro.disconnect();
      window.removeEventListener('orientationchange', measure);
    };
  }, [lines, lineHeightState]);

  return (
    <div
      ref={containerRef}
      className="relative h-full w-full flex flex-col justify-start items-start pt-0 overflow-hidden pr-4"
      style={{ fontFamily: 'Playfair Display, serif', fontWeight: 600 }}
    >
      {/* full-width vertical-only blur overlay below first line (computed left/right) */}
      <div
        className="fixed pointer-events-none z-50"
        style={{
          // position overlay exactly over the panel using viewport coords
          left: overlayAbsLeft != null ? `${overlayAbsLeft}px` : undefined,
          top: overlayAbsTop != null ? `${overlayAbsTop}px` : undefined,
          height: overlayParentHeight != null ? `${overlayParentHeight}px` : undefined,
          width: overlayWidth != null ? `${overlayWidth}px` : undefined,
          // vertical-only gradient; use pixel stop so the gradient starts exactly after the first line
          background: startBlurAt != null && overlayParentTop != null
            ? `linear-gradient(to bottom, rgba(255,255,255,0.0) 0px, rgba(255,255,255,0.0) ${startBlurAt - 1}px, rgba(255,255,255,0.04) ${startBlurAt}px, rgba(255,255,255,0.10) ${Math.min((overlayParentHeight || 0), startBlurAt + 100)}px)`
            : 'linear-gradient(to bottom, rgba(255,255,255,0.0) 0%, rgba(255,255,255,0.04) 8%, rgba(255,255,255,0.10) 35%, rgba(255,255,255,0.22) 100%)',
          WebkitBackdropFilter: 'blur(3px) saturate(1.02) contrast(1.01)',
          backdropFilter: 'blur(3px) saturate(1.02) contrast(1.01)',
          boxShadow: 'inset 0 10px 16px rgba(255,255,255,0.03), inset 0 -8px 12px rgba(0,0,0,0.02)'
        }}
      />

      {lines.map((line, i) => (
        <div
          key={i}
          ref={i === 0 ? firstLineRef : undefined}
          className="flex flex-row gap-x-4 justify-start w-full relative z-10"
          style={{ lineHeight: `${lineHeightState}px` }}
        >
          {line.map((word, idx) => (
            <span
              key={idx}
              className="text-black font-bold tracking-tight select-none"
              style={{ whiteSpace: 'nowrap', fontSize: `${fontSizeState}px` }}
            >
              {word}
            </span>
          ))}
        </div>
      ))}
    </div>
  );
};

export default WordGrid;
