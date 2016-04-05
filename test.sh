input_video_name=$1

output_video_name=${input_video_name}_out.mp4

rm -rf data/${output_video_name}
rm -rf data/${output_video_name}.jpg
rm -rf ../data_out/${input_video_name}.wav
rm -rf ../data_out/${input_video_name}_fin.mp4 

open -W BloggieUnwarpDebug.app --args $input_video_name ../data_out/$output_video_name

ffmpeg -ss 00:00:01 -i data_out/${output_video_name}  -f image2 -y data_out/${input_video_name}_fin.mp4.jpg

ffmpeg -i data/${input_video_name} -map 0:1 data_out/${input_video_name}.wav

ffmpeg -i data_out/${output_video_name} -i data_out/${input_video_name}.wav data_out/${input_video_name}_fin.mp4


rm -rf data_out/${output_video_name}
#rm -rf data/${output_video_name}.jpg
rm -rf data_out/${input_video_name}.wav
