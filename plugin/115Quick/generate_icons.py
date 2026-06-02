from PIL import Image, ImageDraw, ImageFont
import os

def create_icon(size, output_path):
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 绘制圆角矩形背景
    margin = size // 8
    radius = size // 6
    draw.rounded_rectangle(
        [margin, margin, size - margin, size - margin],
        radius=radius,
        fill=(103, 126, 234)  # #667eea
    )
    
    # 绘制渐变效果（简化版）
    for i in range(size // 2):
        alpha = int(255 * (1 - i / (size // 2)))
        color = (118, 75, 162, alpha)  # #764ba2
        draw.line(
            [(margin + i, size - margin - i), (size - margin - i, margin + i)],
            fill=color,
            width=2
        )
    
    # 绘制文字 "115"
    font_size = size // 3
    try:
        font = ImageFont.truetype("arial.ttf", font_size)
    except:
        font = ImageFont.load_default()
    
    text = "115"
    bbox = draw.textbbox((0, 0), text, font=font)
    text_width = bbox[2] - bbox[0]
    text_height = bbox[3] - bbox[1]
    x = (size - text_width) // 2
    y = (size - text_height) // 2 - size // 10
    draw.text((x, y), text, fill="white", font=font)
    
    img.save(output_path)
    print(f"Created {output_path}")

icons_dir = "C:/DEVENV/vueproject/goproject/115Quick_server/plugin/115Quick/public/icons"
os.makedirs(icons_dir, exist_ok=True)

create_icon(16, f"{icons_dir}/icon16.png")
create_icon(48, f"{icons_dir}/icon48.png")
create_icon(128, f"{icons_dir}/icon128.png")
