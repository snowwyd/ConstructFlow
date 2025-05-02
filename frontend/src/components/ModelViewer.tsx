import { Canvas } from '@react-three/fiber';
import { OrbitControls, useGLTF } from '@react-three/drei';

const ModelViewer = ({ url }: { url: string }) => {
    const gltf = useGLTF(url);

    return (
        <Canvas
            camera={{ position: [0, 0, 5], fov: 50 }}
            style={{ width: '100%', height: '100%' }}
        >
            <ambientLight intensity={0.5} />
            <spotLight position={[10, 10, 10]} angle={0.15} penumbra={1} />
            <primitive object={gltf.scene} scale={0.5} />
            <OrbitControls />
        </Canvas>
    );
};

export default ModelViewer;