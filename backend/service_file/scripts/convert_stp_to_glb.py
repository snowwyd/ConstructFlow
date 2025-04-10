# scripts/convert_stp_to_glb.py

import sys
import aspose.cad as cad

def convert_stp_to_glb(input_path, output_path):
    try:
        image = cad.Image.load(input_path)

        options = cad.imageoptions.CadRasterizationOptions()
        options.page_width = 1600.0
        options.page_height = 1600.0

        export_options = cad.imageoptions.GltfOptions()
        export_options.vector_rasterization_options = options

        image.save(output_path, export_options)
        print(output_path)
    except Exception as e:
        print(f"ERROR: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: convert_stp_to_glb.py <input_file> <output_file>", file=sys.stderr)
        sys.exit(1)

    convert_stp_to_glb(sys.argv[1], sys.argv[2])
