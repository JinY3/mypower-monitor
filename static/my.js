// 基于准备好的dom，初始化echarts实例
var myChart = echarts.init(document.getElementById('main'));

// 发送GET请求
fetch('http://157.0.19.2:10063/mypower/data/{{.homeid}}')
    .then(response => response.json())
    .then(data => {
        console.log(data)
        // 指定图表的配置项和数据
        var option = {
            title: {
                text: '{{.homeid}}宿舍剩余电量: ' + data.current + '度'
            },
            tooltip: {},
            legend: {
                data: ['耗电量']
            },
            xAxis: {
                type: 'category',
                data: data.time
            },
            yAxis: {
                type: 'value'
            },
            series: [
                {
                    name: '耗电量',
                    type: 'bar',
                    data: data.value,
                    label: {
                        show: true,
                        position: 'inside'
                    },
                }
            ]
        };

        // 使用刚指定的配置项和数据显示图表。
        myChart.setOption(option);
    })
    .catch(error => {
        console.error('Error:', error);
    });